package memorystore

import (
	"time"

	"github.com/codefluence-x/monorepo/kontext"
)

// MemoryStore is interface to connect to cache infrastructure
type MemoryStore interface {
	// To make it not expire set expiration into 0
	Set(ctx kontext.Context, key string, value []byte, expiration time.Duration) error
	Get(ctx kontext.Context, key string) (Item, error)
}

// Item contain result got from MemoryStore interface
type Item interface {
	Key() string
	Value() []byte
	ExpiresIn() time.Duration
}
