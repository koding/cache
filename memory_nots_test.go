package cache

import "testing"

func TestMemoryCacheNoTSGetSet(t *testing.T) {
	cache := NewMemoryNoTS()
	testCacheGetSet(t, cache)
}

func TestMemoryCacheNoTSSetNX(t *testing.T) {
	cache := NewMemoryNoTS()
	testCacheSetNX(t, cache)
}

func TestMemoryCacheNoTSDelete(t *testing.T) {
	cache := NewMemoryNoTS()
	testCacheDelete(t, cache)
}

func TestMemoryCacheNoTSNilValue(t *testing.T) {
	cache := NewMemoryNoTS()
	testCacheNilValue(t, cache)
}
