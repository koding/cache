package cache

import "testing"

func TestMemoryCacheNoTSGetSet(t *testing.T) {
	cache := NewMemoryNoTS()
	testCacheGetSet(t, cache)
}

func TestMemoryCacheNoTSDelete(t *testing.T) {
	cache := NewMemoryNoTS()
	testCacheDelete(t, cache)
}
