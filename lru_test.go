package cache

import "testing"

func TestLRUCacheGetSet(t *testing.T) {
	cache := NewLRUCache(2)
	testCacheGetSet(t, cache)
}

func TestLRUCacheEviction(t *testing.T) {
	cache := NewLRUCache(2)
	testCacheGetSet(t, cache)

	err := cache.Set("test_key3", "test_data3")
	if err != nil {
		t.Fatal("should not give err while setting item")
	}

	_, err = cache.Get("test_key")
	if err == nil {
		t.Fatal("test_key should not be in the cache")
	}
}

func TestLRUCacheDelete(t *testing.T) {
	cache := NewLRUCache(2)
	testCacheDelete(t, cache)
}
