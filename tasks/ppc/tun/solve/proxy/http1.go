package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"slices"

	"nhooyr.io/websocket"
)

// proxyTCPToHTTP1 waits for a TCP connection on the specified address or uses the one provided,
// and then proxies it to the specified HTTP/1 proxy server over websockets
func proxyTCPToHTTP1(ctx context.Context, inAddr string, tcpConn net.Conn, outAddr string) error {
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

	wsConn, _, err := websocket.Dial(connectCtx, "ws://"+outAddr, &websocket.DialOptions{
		HTTPHeader: http.Header{"User-Agent": []string{}}, // avoid writing any user agent
	})
	if err != nil {
		return fmt.Errorf("dialing target %s over websocket: %w", outAddr, err)
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "bye")

	ctx, cancel := context.WithTimeout(ctx, connTimeout)

	deadline, _ := ctx.Deadline()
	_ = tcpConn.SetDeadline(deadline)

	logger := slog.With("type", "websocket", "in_addr", inAddr, "out_addr", outAddr)
	logger.Info("established websocket proxy connection")

	// Run proxying funcs and wait for them to stop
	proxyFrom := func() error {
		return proxyTCPToWS(ctx, logger, tcpConn, wsConn)
	}

	proxyTo := func() error {
		return proxyWSToTCP(ctx, logger, wsConn, tcpConn)
	}

	err = proxyConcurrently(cancel, tcpConn, proxyFrom, proxyTo)
	logger.Info("ending websocket proxy connection", "reason", err)

	return err
}

func proxyWSToTCP(ctx context.Context, logger *slog.Logger, wsConn *websocket.Conn, tcpConn net.Conn) error {
	for {
		_, r, err := wsConn.Reader(ctx)
		if err != nil {
			return fmt.Errorf("getting next websocket message reader: %w", err)
		}

		buf, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("reading message from websocket: %w", err)
		}

		logger.Info("new message from websocket", "data", buf)

		if _, err := tcpConn.Write(buf); err != nil {
			return fmt.Errorf("writing message to TCP: %w", err)
		}
	}
}

func proxyTCPToWS(ctx context.Context, logger *slog.Logger, tcpConn net.Conn, wsConn *websocket.Conn) error {
	buf := make([]byte, 1<<16)

	for {
		n, readErr := tcpConn.Read(buf)
		if n > 0 {
			b := buf[:n]

			logger.Info("new message from tcp", "data", b)

			if writeErr := wsConn.Write(ctx, websocket.MessageBinary, slices.Clone(b)); writeErr != nil {
				return fmt.Errorf("writing message to websocket: %w", writeErr)
			}
		}

		if readErr != nil {
			return fmt.Errorf("reading message from TCP: %w", readErr)
		}
	}
}
