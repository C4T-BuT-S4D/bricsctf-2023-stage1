package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"slices"
	"time"

	"nhooyr.io/websocket"
)

// http1Protocol wraps a websocket connection, providing proxying on top of an HTTP/1 websocket server
type http1Protocol struct {
	remoteAddr string
	wsConn     *websocket.Conn
}

func (p *http1Protocol) RemoteAddr() string {
	return p.remoteAddr
}

func (p *http1Protocol) Read() ([]byte, error) {
	// context.Background because the connection will be closed by the handler
	_, r, err := p.wsConn.Reader(context.Background())
	if err != nil {
		return nil, err
	}

	return io.ReadAll(r)
}

func (p *http1Protocol) Write(buf []byte) error {
	// context.Background because the connection will be closed by the handler
	return p.wsConn.Write(context.Background(), websocket.MessageBinary, slices.Clone(buf))
}

// http1ProxyServer receives inbound proxy connections over HTTP/1
type http1ProxyServer struct {
	listener   *tcpListener
	proxy      *tcpProxy
	httpServer *http.Server
}

func newHTTP1ProxyServer(listener *tcpListener, forwardProxy *tcpProxy, idleTimeout time.Duration) *http1ProxyServer {
	ps := &http1ProxyServer{
		listener: listener,
		proxy:    forwardProxy,
	}

	ps.httpServer = &http.Server{
		Handler:     http.HandlerFunc(ps.proxyWebsocket),
		IdleTimeout: idleTimeout,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return chainToContext(ctx, c.(*proxyConn))
		},
	}

	return ps
}

// proxyWebsocket is the main handler for HTTP/1 proxy servers.
// It expects a Websocket upgrade request, and then will handle proxying
// using the bidirectional messaging capabilities of websockets.
func (ps *http1ProxyServer) proxyWebsocket(w http.ResponseWriter, r *http.Request) {
	wsConn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close(websocket.StatusNormalClosure, "bye")

	proxyChain := chainFromContext(r.Context())
	if proxyChain == nil {
		// fail fast when connection setup has failed
		http.Error(w, "chain init failed", http.StatusInternalServerError)
		return
	}

	protoConn := &http1Protocol{
		remoteAddr: r.RemoteAddr,
		wsConn:     wsConn,
	}

	if err := ps.proxy.proxy(r.Context(), protoConn, proxyChain); err != nil {
		_ = wsConn.Close(websocket.StatusInternalError, err.Error())
		return
	}
}

func (ps *http1ProxyServer) run() error {
	return ps.httpServer.Serve(ps.listener)
}

func (ps *http1ProxyServer) gracefulStop(ctx context.Context) {
	_ = ps.httpServer.Shutdown(ctx)
}

func (ps *http1ProxyServer) stop() {
	_ = ps.httpServer.Close()
}
