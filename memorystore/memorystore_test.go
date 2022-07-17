package memorystore_test

import (
	"testing"

	"github.com/kodefluence/monorepo/memorystore"
	"github.com/stretchr/testify/assert"
)

func TestMemoryStore(t *testing.T) {

	t.Run("Memcached", func(t *testing.T) {
		memcached := memorystore.FabricateMemcached("main_cache", memorystore.Config{})
		assert.NotNil(t, memcached)

		memcached = memorystore.FabricateMemcached("main_cache", memorystore.Config{})
		assert.NotNil(t, memcached)

		assert.Equal(t, 0, len(memorystore.CloseAll()))
	})
}
