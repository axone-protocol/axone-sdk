package types

import "context"

type Context struct {
	context.Context

	chainID string
}

func (c Context) ChainID() string {
	return c.chainID
}

// WithChainID returns a new context with the given chain ID.
func (c Context) WithChainID(chainID string) Context {
	c.chainID = chainID
	return c
}

// WithContext returns a new context with the given context.
func (c Context) WithContext(ctx context.Context) Context {
	c.Context = ctx
	return c
}
