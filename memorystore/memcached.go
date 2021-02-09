package memorystore

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/codefluence-x/monorepo/exception"
	"github.com/codefluence-x/monorepo/kontext"
)

//go:generate mockgen -source=./memcached.go -destination=./mock/memcached_mock.go -package mock

// Memcached wrap memcache into cache interface
type Memcached struct {
	client MemcachedClient
}

// MemcachedClient wrap default memcache client
type MemcachedClient interface {
	Set(item *memcache.Item) error
	Get(key string) (item *memcache.Item, err error)
}

// AdaptMemcache adapt Cache interface
func AdaptMemcache(client MemcachedClient) *Memcached {
	return &Memcached{client: client}
}

// Set cache value
func (m *Memcached) Set(ktx kontext.Context, key string, value []byte, expiration time.Duration) exception.Exception {
	err := m.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(int(expiration.Seconds())),
	})
	if err != nil {
		return exception.Throw(err)
	}

	return nil
}

// Get cache value
func (m *Memcached) Get(ktx kontext.Context, key string) (Item, exception.Exception) {
	i, err := m.client.Get(key)
	if err == memcache.ErrCacheMiss {
		return nil, exception.Throw(err, exception.WithType(exception.NotFound))
	} else if err != nil {
		return nil, exception.Throw(err)
	}

	return NewCacheItem(i.Key, i.Value, time.Duration(i.Expiration)*time.Second), nil
}
