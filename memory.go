package cache

import "sync"

type MemoryCache struct {
	// Mutex is used for handling the concurrent
	// read/write requests for cache
	sync.Mutex

	// cache holds the cache data
	cache Cache
}

// NewMemoryCache creates an inmemory cache system
// Which everytime will return the true value about a cache hit
func NewMemory() *MemoryCache {
	return &MemoryCache{
		cache: NewMemoryNoTS(),
	}
}

// Get returns a value of a given key if it exists
// and valid for the time being
func (r *MemoryCache) Get(key string) (interface{}, error) {
	r.Lock()
	defer r.Unlock()

	return r.cache.Get(key)
}

// Set will persist a value to the cache or
// override existing one with the new one
func (r *MemoryCache) Set(key string, value interface{}) error {
	r.Lock()
	defer r.Unlock()

	return r.cache.Set(key, value)
}

// Delete deletes a given key if exists
func (r *MemoryCache) Delete(key string) error {
	r.Lock()
	defer r.Unlock()

	return r.cache.Delete(key)
}
