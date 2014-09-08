package cache

import "testing"

func TestMemoryCacheNoTSGetSet(t *testing.T) {
	cache := NewMemoryNoTS()
	cache.Set("test_key", "test_data")
	cache.Set("test_key2", "test_data2")

	data, err := cache.Get("test_key")
	if err != nil {
		t.Fatal("test_key should be in the cache")
	}

	if data != "test_data" {
		t.Fatal("data is not \"test_data\"")
	}

	data, err = cache.Get("test_key2")
	if err != nil {
		t.Fatal("test_key2 should be in the cache")
	}

	if data != "test_data2" {
		t.Fatal("data is not \"test_data2\"")
	}
}

func TestMemoryCacheNoTSDelete(t *testing.T) {
	cache := NewMemoryNoTS()
	cache.Set("test_key", "test_data")
	cache.Set("test_key2", "test_data2")

	err := cache.Delete("test_key3")
	if err != nil {
		t.Fatal("non-exiting item should not give error")
	}

	err = cache.Delete("test_key")
	if err != nil {
		t.Fatal("exiting item should not give error")
	}

	data, err := cache.Get("test_key")
	if err != ErrNotFound {
		t.Fatal("test_key should not be in the cache")
	}

	if data != nil {
		t.Fatal("data should be nil")
	}
}
