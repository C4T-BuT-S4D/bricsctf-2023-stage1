package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// http2Protocol wraps a gRPC stream, providing proxying on top of an HTTP/2 gRPC server
type http2Protocol struct {
	remoteAddr string
	stream     grpc.ServerStream
}

func (p *http2Protocol) RemoteAddr() string {
	return p.remoteAddr
}

func (p *http2Protocol) Read() ([]byte, error) {
	var buf []byte
	err := p.stream.RecvMsg(&buf)
	return buf, err
}

func (p *http2Protocol) Write(buf []byte) error {
	return p.stream.SendMsg(slices.Clone(buf))
}

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

// http2ProxyServer receives inbound proxy connections over HTTP/2
type http2ProxyServer struct {
	listener    *tcpListener
	proxy       *tcpProxy
	http2Server *http2.Server
	grpcServer  *grpc.Server
}

func newHTTP2ProxyServer(listener *tcpListener, forwardProxy *tcpProxy, idleTimeout time.Duration) *http2ProxyServer {
	ps := &http2ProxyServer{
		listener: listener,
		proxy:    forwardProxy,
		http2Server: &http2.Server{
			IdleTimeout: idleTimeout,
		},
	}

	ps.grpcServer = grpc.NewServer(grpc.ForceServerCodec(nopCodec{}))
	ps.grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: "tun.HTTP2",
		HandlerType: (any)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams: []grpc.StreamDesc{{
			StreamName:    "Proxy",
			Handler:       ps.proxyGRPC,
			ServerStreams: true,
			ClientStreams: true,
		}},
		Metadata: nil,
	}, nil)

	return ps
}

// proxyGRPC implements an HTTP/2 proxy handler.
// It exchanges messages using a custom nopCoding which simply returns
// the raw message contained in the data frame.
func (ps *http2ProxyServer) proxyGRPC(_ any, stream grpc.ServerStream) error {
	peer, _ := peer.FromContext(stream.Context())

	proxyChain := chainFromContext(stream.Context())
	if proxyChain == nil {
		// fail fast when connection setup has failed
		return status.Error(codes.Internal, "chain init failed")
	}

	protoConn := &http2Protocol{
		remoteAddr: peer.Addr.String(),
		stream:     stream,
	}

	if err := ps.proxy.proxy(stream.Context(), protoConn, proxyChain); err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

// run manually accepts incoming connections and hands them off to the HTTP/2 server.
// This is used because it's the only way to store the proxy connection in the context.
func (ps *http2ProxyServer) run() error {
	for {
		conn, err := ps.listener.Accept()
		if err != nil {
			return err
		}

		ctx := chainToContext(context.TODO(), conn.(*proxyConn))
		go ps.http2Server.ServeConn(conn, &http2.ServeConnOpts{
			Context: ctx,
			Handler: ps.grpcServer,
		})
	}
}

func (ps *http2ProxyServer) gracefulStop(context.Context) {
	ps.listener.Close() // manually close listener since we're the ones accepting connections
	ps.grpcServer.GracefulStop()
}

func (ps *http2ProxyServer) stop() {
	ps.grpcServer.Stop()
}

func init() {
	encoding.RegisterCodec(nopCodec{})
}
