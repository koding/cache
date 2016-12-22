package cache

import "testing"

func TestLRUNoTSGetSet(t *testing.T) {
	cache := NewLRUNoTS(2)
	testCacheGetSet(t, cache)
}

func TestLRUNoTSSetNX(t *testing.T) {
	cache := NewLRUNoTS(2)
	testCacheSetNX(t, cache)
}

func TestLRUNoTSEviction(t *testing.T) {
	cache := NewLRUNoTS(2)
	testCacheGetSet(t, cache)

	err := cache.Set("test_key3", "test_data3")
	if err != nil {
		t.Fatal("should not give err while setting item")
	}

	_, err = cache.Get("test_key")
	if err == nil {
		t.Fatal("test_key should not be in the cache")
	}

	_, err = cache.SetNX("test_key4", "test_data4")
	if err != nil {
		t.Fatal("should not give err while setting item")
	}

	_, err = cache.Get("test_key2")
	if err == nil {
		t.Fatal("test_key2 should not be in the cache")
	}
}

func TestLRUNoTSDelete(t *testing.T) {
	cache := NewLRUNoTS(2)
	testCacheDelete(t, cache)
}

func TestLRUNoTSNilValue(t *testing.T) {
	cache := NewLRUNoTS(2)
	testCacheNilValue(t, cache)
}
