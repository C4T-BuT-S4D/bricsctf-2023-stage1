package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	dialTimeout = time.Second * 5
	connTimeout = time.Minute * 2
)

// acceptTCPConnection listens on the specified address and accepts a single connection.
// A goroutine is then launched to wait until the context is canceled, at which point the connection gets closed.
func acceptTCPConnection(ctx context.Context, addr string) (net.Conn, context.CancelFunc, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, fmt.Errorf("listening on %s: %w", addr, err)
	}
	defer l.Close()

	type acceptResult struct {
		conn net.Conn
		err  error
	}

	// Accept in separate goroutine in case the context is canceled prematurely
	resultCh := make(chan acceptResult, 1)

	go func() {
		conn, err := l.Accept()
		resultCh <- acceptResult{conn, err}
	}()

	var result acceptResult
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case result = <-resultCh:
	}

	if result.err != nil {
		return nil, nil, fmt.Errorf("accepting connection on %s: %w", addr, result.err)
	}

	cancelTCP := connWithCancel(ctx, result.conn)

	return result.conn, cancelTCP, result.err
}

func connWithCancel(ctx context.Context, conn net.Conn) context.CancelFunc {
	tcpCtx, cancelTcp := context.WithCancel(ctx)

	// Launch another goroutine which will now monitor the context to close the connection
	go func() {
		<-tcpCtx.Done()
		conn.Close()
	}()

	return cancelTcp
}

// proxyConcurrently runs proxying functions in two goroutines,
// and then cancels everything as soon as one of them returns.
func proxyConcurrently(cancel context.CancelFunc, tcpConn net.Conn, from, to func() error) error {
	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer wg.Done()
		errCh <- from()
	}()

	go func() {
		defer wg.Done()
		errCh <- to()
	}()

	err := <-errCh
	cancel()
	_ = tcpConn.SetDeadline(time.Now())
	wg.Wait()

	return err
}
