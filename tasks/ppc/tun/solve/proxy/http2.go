package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"slices"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
)

// nopCodec allows us to use gRPC as HTTP2 without manually managing raw frames
type nopCodec struct{}

func (nopCodec) Name() string {
	return "nop"
}

func (nopCodec) Marshal(v interface{}) ([]byte, error) {
	buf, ok := v.([]byte)
	if !ok {
		return nil, fmt.Errorf("nop: Marshal got unsupported message type %T", v)
	} else if buf == nil {
		return nil, errors.New("nop: Marshal doesn't handle nil messages")
	}

	return buf, nil
}

func (nopCodec) Unmarshal(data []byte, v interface{}) error {
	buf, ok := v.(*[]byte)
	if !ok {
		return fmt.Errorf("nop: Unmarshal got unsupported message type %T", v)
	} else if buf == nil {
		return errors.New("nop: Unmarshal doesn't handle nil messages")
	}

	*buf = data
	return nil
}

// proxyTCPToHTTP2 waits for a TCP connection on the specified address or uses the one provided,
// and then proxies it to the specified HTTP/2 proxy server over gRPC
func proxyTCPToHTTP2(ctx context.Context, inAddr string, tcpConn net.Conn, outAddr string) error {
	var cancelTCP context.CancelFunc
	if tcpConn == nil {
		var err error
		tcpConn, cancelTCP, err = acceptTCPConnection(ctx, inAddr)
		if err != nil {
			return err
		}
	} else {
		cancelTCP = connWithCancel(ctx, tcpConn)
	}

	defer cancelTCP()

	// Establish connection to the websocket proxy
	connectCtx, cancelConnect := context.WithTimeout(ctx, dialTimeout)
	defer cancelConnect()

	cc, err := grpc.DialContext(connectCtx, outAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype("nop")),
	)
	if err != nil {
		return fmt.Errorf("dialing target %s over gRPC: %w", outAddr, err)
	}
	defer cc.Close()

	ctx, cancel := context.WithTimeout(ctx, connTimeout)

	proxyStream, err := cc.NewStream(ctx, &grpc.StreamDesc{ServerStreams: true, ClientStreams: true}, "/tun.HTTP2/Proxy")
	if err != nil {
		cancel()
		return fmt.Errorf("initiating gRPC proxy stream: %w", err)
	}

	deadline, _ := ctx.Deadline()
	_ = tcpConn.SetDeadline(deadline)

	logger := slog.With("type", "grpc", "in_addr", inAddr, "out_addr", outAddr)
	logger.Info("established gRPC proxy connection")

	// Run proxying funcs and wait for them to stop
	proxyFrom := func() error {
		return proxyTCPToGRPC(ctx, logger, tcpConn, proxyStream)
	}

	proxyTo := func() error {
		return proxyGRPCToTCP(ctx, logger, proxyStream, tcpConn)
	}

	err = proxyConcurrently(cancel, tcpConn, proxyFrom, proxyTo)
	logger.Info("ending gRPC proxy connection", "reason", err)

	return err
}

func proxyGRPCToTCP(ctx context.Context, logger *slog.Logger, proxyStream grpc.ClientStream, tcpConn net.Conn) error {
	var buf []byte

	for {
		if err := proxyStream.RecvMsg(&buf); err != nil {
			return fmt.Errorf("reading message from stream: %w", err)
		}

		logger.Info("new message from gRPC", "data", buf)

		if _, err := tcpConn.Write(buf); err != nil {
			return fmt.Errorf("writing message to TCP: %w", err)
		}
	}
}

func proxyTCPToGRPC(ctx context.Context, logger *slog.Logger, tcpConn net.Conn, proxyStream grpc.ClientStream) error {
	buf := make([]byte, 1<<16)

	for {
		n, readErr := tcpConn.Read(buf)
		if n > 0 {
			b := buf[:n]

			logger.Info("new message from tcp", "data", b)

			if writeErr := proxyStream.SendMsg(slices.Clone(b)); writeErr != nil {
				return fmt.Errorf("writing message to stream: %w", writeErr)
			}
		}

		if readErr != nil {
			return fmt.Errorf("reading message from TCP: %w", readErr)
		}
	}
}

func init() {
	encoding.RegisterCodec(nopCodec{})
}
