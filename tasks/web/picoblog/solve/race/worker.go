package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"log/slog"
)

var httpEnding = []byte("\r\n\r\n")

type worker struct {
	name         string
	wg           *sync.WaitGroup
	phase1Notify chan struct{}
	errNotify    chan<- error
	conn         net.Conn
	init         []byte
	last         []byte
}

func newWorker(name string, addr, host string, wg *sync.WaitGroup, phase1Notify chan struct{}, errNotify chan<- error, req string) (*worker, error) {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return nil, fmt.Errorf("dialing target: %w", err)
	}

	init, last := []byte(req[:len(req)-1]), req[len(req)-1]

	return &worker{
		name:         name,
		wg:           wg,
		phase1Notify: phase1Notify,
		errNotify:    errNotify,
		conn:         conn,
		init:         init,
		last:         []byte{last},
	}, nil
}

func (w *worker) run() {
	resp, err := w.runAttack()
	if err != nil {
		slog.Warn("attack failed", "worker", w.name, "error", err)
		w.errNotify <- fmt.Errorf("worker %s: %w", w.name, err)
		return
	}

	slog.Info("attack success", "worker", w.name, "resp", string(resp))

	_ = w.conn.Close()

	// notify that the worker was succcessful
	w.wg.Done()
}

func (w *worker) runAttack() ([]byte, error) {
	if _, err := w.conn.Write(w.init); err != nil {
		return nil, fmt.Errorf("writing init part: %w", err)
	}

	// notify coordinator that first part is ready
	w.wg.Done()

	// wait for signal that all workers are ready
	<-w.phase1Notify

	if _, err := w.conn.Write(w.last); err != nil {
		return nil, fmt.Errorf("writing last part: %w", err)
	}

	// Read response after all request bytes have been sent
	var resp []byte
	r := bufio.NewReader(w.conn)

	for i := 0; !bytes.HasSuffix(resp, httpEnding); i++ {
		line, err := r.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("reading response line %d: %w", i, err)
		}

		resp = append(resp, line...)
	}

	// a bit of the response body
	extra := make([]byte, 50)
	_, _ = r.Read(extra)

	resp = append(resp, extra...)

	return resp, nil
}
