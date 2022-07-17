package memorystore

import (
	"time"

	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
)

// MemoryStore is interface to connect to cache infrastructure
type MemoryStore interface {
	// To make it not expire set expiration into 0
	Set(ktx kontext.Context, key string, value []byte, expiration time.Duration) exception.Exception
	Get(ktx kontext.Context, key string) (Item, exception.Exception)
}

// Item contain result got from MemoryStore interface
type Item interface {
	Key() string
	Value() []byte
	ExpiresIn() time.Duration
}
