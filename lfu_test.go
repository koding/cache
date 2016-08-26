package cache

import (
	"fmt"
	"testing"
)

func TestLFUNGetSet(t *testing.T) {
	cache := NewLFU(2)
	testCacheGetSet(t, cache)
}

func TestLFUDelete(t *testing.T) {
	cache := NewLFU(2)
	testCacheDelete(t, cache)
}

func TestLFUNilValue(t *testing.T) {
	cache := NewLFU(2)
	testCacheNilValue(t, cache)
}

func TestLFUEviction(t *testing.T) {
	cache := NewLFU(2)
	testCacheGetSet(t, cache)

	_, err := cache.Get("test_key2")
	if err != nil {
		t.Fatal("test_key2 should be in the cache")
	}
	// get-> test_key should not be in cache after insert test_key3
	err = cache.Set("test_key3", "test_data3")
	if err != nil {
		t.Fatal("should not give err while setting item")
	}

	_, err = cache.Get("test_key")
	if err == nil {
		t.Fatal("test_key should not be in the cache")
	}
}

//
// BENCHMARK OPERATIONS
//

func BenchmarkLFUSet1000(b *testing.B) {
	cache := NewLFU(5)
	for n := 0; n < b.N; n++ {
		testSetNTimes(cache, 1000)
	}
}
func BenchmarkLFUGet1000(b *testing.B) {
	cache := NewLFU(5)
	for n := 0; n < b.N; n++ {
		testGetNTimes(cache, 1000)
	}
}
func BenchmarkLFUSetDelete1000(b *testing.B) {
	cache := NewLFU(5)
	for n := 0; n < b.N; n++ {
		testSetDeleteNTimes(cache, 1000)
	}
}
func BenchmarkLFUSetGet1000(b *testing.B) {
	cache := NewLFU(5)
	for n := 0; n < b.N; n++ {
		testSetGetNTimes(cache, 1000)
	}
}

func testSetNTimes(cache Cache, n int) {
	for i := 0; i < n; i++ {
		cache.Set("keyBench", i)
	}
}

func testGetNTimes(cache Cache, n int) {
	_, err := cache.Get("keyBench")
	if err != nil && err != ErrNotFound {
		fmt.Println("Occurred error while getting from cache")
		return
	}
	if err == ErrNotFound {
		cache.Set("keyBench", "test")
	}

	for i := 0; i < n-1; i++ {
		cache.Get("keyBench")
	}
}

func testSetDeleteNTimes(cache Cache, n int) {
	for i := 0; i < n; i++ {
		cache.Set("keyBench", i)
		cache.Delete("keyBench")
	}
}
func testSetGetNTimes(cache Cache, n int) {
	for i := 0; i < n; i++ {
		cache.Set("keyBench", i)
		cache.Get("keyBench")
	}
}
