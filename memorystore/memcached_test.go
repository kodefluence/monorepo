package memorystore_test

import (
	"errors"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/golang/mock/gomock"
	"github.com/kodefluence/monorepo/exception"
	"github.com/kodefluence/monorepo/kontext"
	"github.com/kodefluence/monorepo/memorystore"
	"github.com/kodefluence/monorepo/memorystore/mock"
	"github.com/stretchr/testify/assert"
)

func TestMemcached(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	t.Run("Set", func(t *testing.T) {
		t.Run("When connection to memcached complete it will return nil", func(t *testing.T) {
			memcacheClient := mock.NewMockMemcachedClient(mockCtrl)
			memcacheClient.EXPECT().Set(gomock.Any()).Return(nil)

			memcachedAdapter := memorystore.AdaptMemcache(memcacheClient)
			assert.Nil(t, memcachedAdapter.Set(kontext.Fabricate(), "key", []byte("value"), 0))
		})

		t.Run("When connection to memcached error it will return exception", func(t *testing.T) {
			memcacheClient := mock.NewMockMemcachedClient(mockCtrl)
			memcacheClient.EXPECT().Set(gomock.Any()).Return(errors.New("unexpected error"))

			memcachedAdapter := memorystore.AdaptMemcache(memcacheClient)
			assert.NotNil(t, memcachedAdapter.Set(kontext.Fabricate(), "key", []byte("value"), 0))
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("It return CacheItem", func(t *testing.T) {
			key := "key"
			value := []byte("value")
			expiration := time.Second

			memcacheClient := mock.NewMockMemcachedClient(mockCtrl)
			memcacheClient.EXPECT().Get(key).Return(&memcache.Item{
				Key:        key,
				Value:      value,
				Expiration: int32(int(expiration.Seconds())),
			}, nil)

			memcachedAdapter := memorystore.AdaptMemcache(memcacheClient)

			cacheItem, err := memcachedAdapter.Get(kontext.Fabricate(), key)
			assert.Nil(t, err)
			assert.Equal(t, key, cacheItem.Key())
			assert.Equal(t, value, cacheItem.Value())
			assert.Equal(t, expiration, cacheItem.ExpiresIn())
		})

		t.Run("When there is unexpected error then it will return the error", func(t *testing.T) {
			memcacheClient := mock.NewMockMemcachedClient(mockCtrl)
			memcacheClient.EXPECT().Get("key").Return(nil, errors.New("unexpected error"))

			memcachedAdapter := memorystore.AdaptMemcache(memcacheClient)

			cacheItem, err := memcachedAdapter.Get(kontext.Fabricate(), "key")
			assert.Nil(t, cacheItem)
			assert.NotNil(t, err)
			assert.Equal(t, exception.Unexpected, err.Type())
		})

		t.Run("When it's cache not found error then it will return not found exception", func(t *testing.T) {
			memcacheClient := mock.NewMockMemcachedClient(mockCtrl)
			memcacheClient.EXPECT().Get("key").Return(nil, memcache.ErrCacheMiss)

			memcachedAdapter := memorystore.AdaptMemcache(memcacheClient)

			cacheItem, err := memcachedAdapter.Get(kontext.Fabricate(), "key")
			assert.Nil(t, cacheItem)
			assert.NotNil(t, err)
			assert.Equal(t, exception.NotFound, err.Type())
		})
	})
}
