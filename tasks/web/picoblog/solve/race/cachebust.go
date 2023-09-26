package main

import (
	"crypto/tls"
	"fmt"
)

func cacheBust(addr string, host string, cacheBustReq string) error {
	bustConn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Errorf("dialing for cachebust: %w", err)
	}

	if _, err := bustConn.Write([]byte(cacheBustReq)); err != nil {
		return fmt.Errorf("writing cachebust request: %w", err)
	}

	cacheBustResp := make([]byte, 200)
	if _, err := bustConn.Read(cacheBustResp); err != nil {
		return fmt.Errorf("reading cachebust response: %w", err)
	}

	_ = bustConn.Close()

	return nil
}
