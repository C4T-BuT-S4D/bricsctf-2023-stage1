package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"sync"

	"log/slog"
)

const (
	portMax = 65535
	portMin = 1000
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(os.Args[1], os.Args[2]); err != nil {
		slog.Error("encountered fatal error", "error", err)
		os.Exit(1)
	}
}

func run(inAddr, outAddr string) error {
	inListener, err := net.Listen("tcp", inAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", inAddr, err)
	}
	defer inListener.Close()

	for {
		inConn, err := inListener.Accept()
		if err != nil {
			return fmt.Errorf("accepting new connection on %s: %w", inAddr, err)
		}

		go func() {
			err := runProxy(inAddr, inConn, outAddr)
			slog.Error(fmt.Sprintf("proxy from %s ended", inConn.RemoteAddr().String()), "reason", err)
		}()
	}
}

// runProxy runs the actual proxying process for an inbound TCP connection
func runProxy(inAddr string, inConn net.Conn, outAddr string) error {
	conn, err := net.Dial("tcp", outAddr)
	if err != nil {
		return err
	}

	chain, err := readChainOrder(conn)
	if err != nil {
		return err
	}

	slog.Info("read proxy chain config", "chain", chain)

	// Build proxies in reverse order
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(len(chain) + 1)

	errCh := make(chan error, len(chain)+1)

	// Proxies will form chain starting from inAddr and ending at a random lastOut address,
	// from which then the fully encapsulated contents will be proxied to the actual proxy connection
	previousOut := inAddr
	lastOut := randLocalAddr()

	slices.Reverse(chain)
	for i, proxyType := range chain {
		var proxyFunc func(ctx context.Context, inAddr string, tcpConn net.Conn, outAddr string) error

		if proxyType == "http1" {
			proxyFunc = proxyTCPToHTTP1
		} else if proxyType == "http2" {
			proxyFunc = proxyTCPToHTTP2
		} else {
			return fmt.Errorf("unsupported proxy type %q encountered in chain", proxyType)
		}

		nextAddr := randLocalAddr()
		if i == len(chain)-1 {
			nextAddr = lastOut
		}

		go func(inAddr string, tcpConn net.Conn, outAddr string) {
			defer wg.Done()
			errCh <- proxyFunc(ctx, inAddr, tcpConn, outAddr)
		}(previousOut, inConn, nextAddr)

		previousOut = nextAddr
		inConn = nil
	}

	// Run additional listener which will simply send the final packets to the proxy
	go func() {
		defer wg.Done()
		errCh <- copyToProxy(ctx, lastOut, conn)
	}()

	wg.Wait()
	err = <-errCh
	slog.Info("shutting everything down", "reason", err)

	return nil
}

// readChainOrder reads the proxy chain config from the first line received from the connection.
func readChainOrder(conn net.Conn) ([]string, error) {
	var buf []byte

	b := make([]byte, 1)
	for len(buf) == 0 || buf[len(buf)-1] != '\n' {
		_, err := conn.Read(b)
		if err != nil {
			return nil, fmt.Errorf("reading proxy chain config preface: %w", err)
		}

		buf = append(buf, b[0])
	}

	var chainConfig struct {
		Order []string `json:"order"`
	}

	if err := json.Unmarshal(buf, &chainConfig); err != nil {
		return nil, fmt.Errorf("unmarshaling proxy chain config: %w", err)
	}

	return chainConfig.Order, nil
}

// copyToProxy accepts a single connection on inAddr and then copies everything
// between the received connection and the proxy connection.
func copyToProxy(ctx context.Context, inAddr string, proxyConn net.Conn) error {
	inConn, cancelTCP, err := acceptTCPConnection(ctx, inAddr)
	if err != nil {
		return err
	}

	slog.Info("copying final packets", "from", inAddr, "to", proxyConn.RemoteAddr().String())

	proxyFrom := func() error {
		_, err := io.Copy(proxyConn, inConn)
		return err
	}

	proxyTo := func() error {
		_, err := io.Copy(inConn, proxyConn)
		return err
	}

	return proxyConcurrently(cancelTCP, proxyConn, proxyFrom, proxyTo)
}

func randLocalAddr() string {
	return "127.0.0.1:" + strconv.Itoa(rand.Intn(portMax-portMin+1)+portMin)
}
