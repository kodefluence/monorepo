package memorystore

import "time"

// CacheItem wrap cache return object
type CacheItem struct {
	key        string
	value      []byte
	expiration time.Duration
}

// NewCacheItem return new cache item object
func NewCacheItem(key string, value []byte, expiration time.Duration) *CacheItem {
	return &CacheItem{key: key, value: value, expiration: expiration}
}

// Key of cache object
func (c *CacheItem) Key() string {
	return c.key
}

// Value of cache object
func (c *CacheItem) Value() []byte {
	return c.value
}

// ExpiresIn contain expired information of cache object
func (c *CacheItem) ExpiresIn() time.Duration {
	return c.expiration
}
