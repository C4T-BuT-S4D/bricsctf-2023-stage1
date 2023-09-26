package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"time"

	"gitlab.com/NebulousLabs/fastrand"
)

// edgeConfig specifies the internal proxy and target addresses to be used by the edge proxy.
type edgeConfig struct {
	HTTP1ProxyAddr string `json:"http1_proxy_addr"`
	HTTP2ProxyAddr string `json:"http2_proxy_addr"`
	FinalAddr      string `json:"final_addr"`
	ChainLength    int    `json:"chain_length"`
}

// tcpListener is a specialized net.Listener which is able to accept chain proxy connections.
type tcpListener struct {
	netListener   net.Listener
	acceptTimeout time.Duration
	edgeType      string
	edgeCfg       *edgeConfig
}

// proxyConn wraps net.proxyConn with an additional method to retrieve the proxy chain for this connection.
type proxyConn struct {
	net.Conn
	chainAddrs []string
}

// chain returns the proxy address chain.
func (c *proxyConn) chain() []string {
	return c.chainAddrs
}

// listenEdgeTCP listens on the specified address and configures a Listener to be used as an edge proxy.
func listenEdgeTCP(addr string, edgeType string, config edgeConfig) (*tcpListener, error) {
	listener, err := listenInternalTCP(addr, 0)
	if err != nil {
		return nil, err
	}

	listener.edgeType = edgeType
	listener.edgeCfg = &config

	return listener, nil
}

// listenInternalTCP listens on the specified address and configures a Listener to be used as an internal proxy.
func listenInternalTCP(addr string, acceptTimeout time.Duration) (*tcpListener, error) {
	netListener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("listening on %q: %w", addr, err)
	}

	return &tcpListener{
		netListener:   netListener,
		acceptTimeout: acceptTimeout,
	}, nil
}

// Accept accepts a new connection using the underlying net.Listener, and then,
// depending on whether this is an edge proxy listener or an internal proxy listener,
// generates a proxy chain and sends it to the client, or reads the proxy chain from the connection.
// *Conn is returned as the net.Conn implementation and can be used to retrieve the proxy chain.
func (l *tcpListener) Accept() (net.Conn, error) {
	conn, err := l.netListener.Accept()
	if err != nil {
		return nil, err
	}

	var chainAddrs []string

	if l.edgeCfg != nil {
		chainAddrs, err = l.writeChainConfig(conn)

		if err != nil {
			slog.Error(
				"edge proxy failed to write chain configuration to new connection",
				"remote_addr", conn.RemoteAddr().String(),
				"error", err,
			)
		}
	} else {
		chainAddrs, err = l.readChainConfig(conn)

		if err != nil {
			slog.Error(
				"internal proxy failed to read chain configuration from new connection",
				"remote_addr", conn.RemoteAddr().String(),
				"error", err,
			)
		}
	}

	// here chainAddrs can be nil if initialization failed,
	// but we shouldn't fail the whole Accept because this can be a temporary error
	return &proxyConn{Conn: conn, chainAddrs: chainAddrs}, nil
}

// writeChainConfig generates a new chain config, writes the proxy types to the connection,
// and returns the generated addresses.
func (l *tcpListener) writeChainConfig(conn net.Conn) ([]string, error) {
	// chainConfig is needed to make it easier to understand what the list means.
	// This is a CTF task, after all.
	type chainConfig struct {
		Order []string `json:"order"`
	}

	// generate chain and send it to the client
	chainTypes, chainAddrs := l.generateChain()

	// only types are sent to the client
	chainData, err := json.Marshal(chainConfig{Order: chainTypes})
	if err != nil {
		return nil, fmt.Errorf("marshaling chain config for client: %w", err)
	}

	// don't return error from write even if it occurs to avoid breaking servers due to uncooperative clients
	if _, err := conn.Write(append(chainData, '\n')); err != nil {
		return nil, fmt.Errorf("writing chain config to client: %w", err)
	}

	return chainAddrs, nil
}

// readChainConfig reads the chain config from the connection and returns the next addresses to use.
func (l *tcpListener) readChainConfig(conn net.Conn) ([]string, error) {
	// read chain from connection, the internal proxies receive a list of next addresses
	var chainData []byte

	// safeguard timeout to avoid breaking this internal proxy
	// if the connection from the previous proxy is corrupt
	_ = conn.SetReadDeadline(time.Now().Add(l.acceptTimeout))

	// read 1-by-1 because net.Conn can't put back read bytes
	buf := make([]byte, 1)

	for len(chainData) == 0 || chainData[len(chainData)-1] != '\n' {
		n, err := conn.Read(buf)
		if n > 0 {
			chainData = append(chainData, buf...)
		}

		if err != nil {
			return nil, fmt.Errorf("reading chain config from conn: %w", err)
		}
	}

	// reset deadline after reading the configuration
	_ = conn.SetDeadline(time.Time{})

	var chainAddrs []string

	if err := json.Unmarshal(chainData, &chainAddrs); err != nil {
		return nil, fmt.Errorf("unmarshaling chain config (%s): %w", base64.RawStdEncoding.EncodeToString(chainData), err)
	}

	return chainAddrs, nil
}

func (l *tcpListener) generateChain() (types, addrs []string) {
	chainTypes := make([]string, l.edgeCfg.ChainLength+1)
	chainAddrs := make([]string, l.edgeCfg.ChainLength+1)

	chainTypes[0] = l.edgeType
	chainAddrs[l.edgeCfg.ChainLength] = l.edgeCfg.FinalAddr

	// Always fill the chain with half of each proxy type, then shuffle them,
	// instead of simply filling with random types, so that all chains are "fair".
	for i := 0; i < l.edgeCfg.ChainLength; i++ {
		if i < l.edgeCfg.ChainLength/2 {
			chainTypes[i+1] = "http1"
			chainAddrs[i] = l.edgeCfg.HTTP1ProxyAddr
		} else {
			chainTypes[i+1] = "http2"
			chainAddrs[i] = l.edgeCfg.HTTP2ProxyAddr
		}
	}

	fastrand.Shuffle(l.edgeCfg.ChainLength, func(i, j int) {
		chainAddrs[i], chainAddrs[j] = chainAddrs[j], chainAddrs[i]
		chainTypes[i+1], chainTypes[j+1] = chainTypes[j+1], chainTypes[i+1]
	})

	return chainTypes, chainAddrs
}

// Close closes the underlying net.Listener.
func (l *tcpListener) Close() error {
	return l.netListener.Close()
}

// Addr returns the underlying net.Listener's Addr.
func (l *tcpListener) Addr() net.Addr {
	return l.netListener.Addr()
}
