package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"
)

// protocolConn describes an inbound protocol which supports reading and writing messages in a streaming manner.
type protocolConn interface {
	RemoteAddr() string
	Read() ([]byte, error)
	Write([]byte) error
}

// tcpProxy proxies incoming protocol traffic to TCP.
type tcpProxy struct {
	bufSize     int
	initTimeout time.Duration
	connTimeout time.Duration
}

// newTCPProxy initializes a new forward TCP proxy which will use the specified timeouts for dialing and connection handling.
func newTCPProxy(bufSize int, initTimeout, connTimeout time.Duration) *tcpProxy {
	return &tcpProxy{
		bufSize:     bufSize,
		initTimeout: initTimeout,
		connTimeout: connTimeout,
	}
}

// proxy establishes a connection to the next address in the chain using TCP and then proxies the inbound protocol connection.
func (p *tcpProxy) proxy(ctx context.Context, protoConn protocolConn, chain []string) error {
	logger := slog.With("remote_addr", protoConn.RemoteAddr())

	if len(chain) < 1 {
		logger.Error("receieved empty proxy chain")
		return errors.New("bad proxy chain")
	}

	// targetAddr is the next proxy or the final addr
	targetAddr := chain[0]
	chain = chain[1:]

	logger = logger.With("target_addr", targetAddr)
	logger.Info("initiating new proxy connection")

	tcpConn, err := net.DialTimeout("tcp", targetAddr, p.initTimeout)
	if err != nil {
		logger.Error("failed to dial target", "error", err)
		return errors.New("connection error")
	}
	defer tcpConn.Close()

	// if chain isn't empty, then we should transmit it to the next proxy
	if len(chain) != 0 {
		if err := p.writeChain(tcpConn, chain); err != nil {
			logger.Error("failed to initialize connection to next proxy", "error", err)
			return errors.New("proxy init error")
		}
	}

	logger = logger.With("local_addr", tcpConn.LocalAddr().String())

	ctx, cancel := context.WithTimeout(ctx, p.connTimeout)
	defer cancel()

	deadline, _ := ctx.Deadline()
	_ = tcpConn.SetDeadline(deadline)

	// Proxy connections both ways. No waitgroup is used because the proto connections will be closed by their handlers,
	// i.e. wsConn.Close() called in case of websockets, gRPC stream handler returned in case of gRPC
	notifyCh := make(chan error, 2)

	go func() {
		err := fmt.Errorf("proto to TCP: %w", p.protoToTCP(protoConn, tcpConn))
		notifyCh <- err
		logger.Info("proto to TCP goroutine exiting", "reason", err)
	}()

	go func() {
		err := fmt.Errorf("TCP to proto: %w", p.tcpToProto(tcpConn, protoConn))
		notifyCh <- err
		logger.Info("TCP to proto goroutine exiting", "reason", err)
	}()

	// Once at least one of the ends is done, cancel the context to close the other.
	// Select on ctx.Done() as well to avoid any potential blocks in the proxy goroutines.
	select {
	case err = <-notifyCh:
	case <-ctx.Done():
		err = ctx.Err()
	}

	logger.Info("closing proxy connection", "reason", err)

	return nil
}

// writeChain sends the proxy chain info over the established connection to the next proxy.
func (p *tcpProxy) writeChain(tcpConn net.Conn, chain []string) error {
	chainData, err := json.Marshal(chain)
	if err != nil {
		return fmt.Errorf("marshaling chain config: %w", err)
	}

	_ = tcpConn.SetWriteDeadline(time.Now().Add(p.initTimeout))

	if _, err := tcpConn.Write(append(chainData, '\n')); err != nil {
		return fmt.Errorf("writing chain config to next proxy: %w", err)
	}

	_ = tcpConn.SetDeadline(time.Time{})

	return nil
}

func (p *tcpProxy) protoToTCP(protoConn protocolConn, tcpConn net.Conn) error {
	for {
		buf, err := protoConn.Read()
		if err != nil {
			return err
		}

		if _, err := tcpConn.Write(buf); err != nil {
			return err
		}
	}
}

func (p *tcpProxy) tcpToProto(tcpConn net.Conn, protoConn protocolConn) error {
	buf := make([]byte, p.bufSize)

	for {
		n, readErr := tcpConn.Read(buf)
		if n > 0 {
			if writeErr := protoConn.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}

		if readErr != nil {
			return readErr
		}
	}
}
