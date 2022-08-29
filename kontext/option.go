package kontext

import (
	"context"
)

// A Config for Context creation
type Config struct {
	ctx context.Context
}

// Option when fabricating Context
type Option func(*Config)

// WithDefaultContext give kontext option to wrap current context with default context
func WithDefaultContext(ctx context.Context) Option {
	return func(c *Config) {
		c.ctx = ctx
	}
}
