package memorystore_test

import (
	"testing"
	"time"

	"github.com/kodefluence/monorepo/memorystore"
	"github.com/stretchr/testify/assert"
)

func TestCacheItem(t *testing.T) {
	key := "some-key"
	value := []byte("some value")
	expiration := time.Duration(0)

	cacheItem := memorystore.NewCacheItem(key, value, expiration)

	t.Run("Key", func(t *testing.T) {
		t.Run("Return cache key", func(t *testing.T) {
			assert.Equal(t, key, cacheItem.Key())
		})
	})

	t.Run("Value", func(t *testing.T) {
		t.Run("Return cache value", func(t *testing.T) {
			assert.Equal(t, value, cacheItem.Value())
		})
	})

	t.Run("ExpiresIn", func(t *testing.T) {
		t.Run("Return cache expiration information", func(t *testing.T) {
			assert.Equal(t, expiration, cacheItem.ExpiresIn())
		})
	})
}
