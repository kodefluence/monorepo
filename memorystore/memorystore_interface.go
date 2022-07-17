package memorystore

import (
	"fmt"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/kodefluence/monorepo/exception"
)

var instanceList = &sync.Map{}

// FabricateMemcached will fabricate memcached and wrap it into MemoryStore interface
func FabricateMemcached(instanceName string, config Config) MemoryStore {
	if val, ok := instanceList.Load(fmt.Sprintf("memcached-%s", instanceName)); ok {
		return AdaptMemcache(val.(*memcache.Client))
	}

	client := memcache.New(fmt.Sprintf("%s:%s", config.Host, config.Port))

	instanceList.Store(fmt.Sprintf("memcached-%s", instanceName), client)

	return AdaptMemcache(client)
}

// CloseAll MemoryStore connection
func CloseAll() []exception.Exception {
	var excs []exception.Exception

	return excs
}
