package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"
)

const (
	edgeConfigFile    = "edge.json"
	proxyBufSize      = 1 << 14         // 16kb buffer (because websockets can send maximum 32kb) for reading from TCP conn
	proxyInitTimeout  = time.Second * 5 // timeout for proxy connection initialization (accept, dial, etc)
	proxyConnTimeout  = time.Minute     // total timeout for proxy connection after initialization
	serverIdleTimeout = time.Minute * 2
	shutdownTimeout   = time.Second * 10
)

type proxyServer interface {
	run() error
	gracefulStop(ctx context.Context)
	stop()
}

func main() {
	proto := flag.String("proto", "", "inbound proxy protocol (http1 or http2)")
	listen := flag.String("listen", "", "address to listen on")
	mode := flag.String("mode", "", "proxy mode (edge or internal)")
	flag.Parse()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(*proto, *listen, *mode); err != nil {
		slog.Error("encountered fatal error", "error", err)
		os.Exit(1)
	}
}

func run(proto, listen, mode string) error {
	if proto != "http1" && proto != "http2" {
		return fmt.Errorf("unsupported proto %q specified", proto)
	} else if mode != "edge" && mode != "internal" {
		return fmt.Errorf("unsupported mode %q specified", mode)
	}

	listener, err := createListener(proto, listen, mode)
	if err != nil {
		return err
	}

	slog.Info("started proxy listener", "listen_addr", listen, "mode", mode)

	server := createServer(listener, proto)

	slog.Info("running proxy server", "proto", proto)

	// Run the actual server, accepting the incoming connections and proxying them
	errCh := make(chan error, 1)
	go func() {
		if err := server.run(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Wait for graceful shutdown or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	var runErr error
	select {
	case <-shutdown:
		slog.Info("performing graceful shutdown")
	case runErr = <-errCh:
		slog.Error("unexpected server error, performing shutdown", "error", runErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Handle stop in separate goroutine to avoid blocking
	done := make(chan struct{})
	go func() {
		defer close(done)
		server.gracefulStop(ctx)
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}

	server.stop()

	return runErr
}

func createListener(proto, listen, mode string) (*tcpListener, error) {
	if mode == "internal" {
		return listenInternalTCP(listen, proxyInitTimeout)
	}

	configData, err := os.ReadFile(edgeConfigFile)
	if err != nil {
		return nil, fmt.Errorf("reading edge config file: %w", err)
	}

	var edgeConfig edgeConfig

	if err := json.Unmarshal(configData, &edgeConfig); err != nil {
		return nil, fmt.Errorf("unmarshaling edge config: %w", err)
	}

	return listenEdgeTCP(listen, proto, edgeConfig)
}

func createServer(listener *tcpListener, proto string) proxyServer {
	proxy := newTCPProxy(proxyBufSize, proxyInitTimeout, proxyConnTimeout)

	if proto == "http1" {
		return newHTTP1ProxyServer(listener, proxy, serverIdleTimeout)
	}

	return newHTTP2ProxyServer(listener, proxy, serverIdleTimeout)
}
