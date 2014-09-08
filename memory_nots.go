package cache

type MemoryCacheNoTS struct {
	// items holds the cache data
	items map[string]interface{}
}

// NewMemoryNoTS creates an inmemory cache system
func NewMemoryNoTS() *MemoryCacheNoTS {
	return &MemoryCacheNoTS{
		items: map[string]interface{}{},
	}
}

// Get returns a value of a given key if it exists
// and valid for the time being
func (r *MemoryCacheNoTS) Get(key string) (interface{}, error) {
	value, ok := r.items[key]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}

// Set will persist a value to the cache or
// override existing one with the new one
func (r *MemoryCacheNoTS) Set(key string, value interface{}) error {
	r.items[key] = value
	return nil
}

// Delete deletes a given key, it doesnt return error if the item is not in the
// system
func (r *MemoryCacheNoTS) Delete(key string) error {
	delete(r.items, key)
	return nil
}
