package main

import (
	"context"
)

type chainCtxKey struct{}

// chainToContext saves the proxy.Conn proxy chain into the context.
func chainToContext(ctx context.Context, conn *proxyConn) context.Context {
	chain := conn.chain()
	return context.WithValue(ctx, chainCtxKey{}, chain)
}

// chainFromContext retrieves the saved proxy chain from the context.
func chainFromContext(ctx context.Context) []string {
	return ctx.Value(chainCtxKey{}).([]string)
}
