package cache

import (
	"container/list"
	"fmt"
)

type LFUNoTS struct {
	// list holds all items in a linked list
	frequencyList *list.List

	// holds the all cache values
	cache Cache

	// size holds the limit of the LFU cache
	size int

	currentSize int
}

type cacheItem struct {
	// key of cache value
	k string

	// value of cache value
	v interface{}

	// holds the frequency elements
	freqElement *list.Element
}

// NewLRUNoTS creates a new LFU cache struct for further cache operations. Size
// is used for limiting the upper bound of the cache
func NewLFUNoTS(size int) Cache {
	if size < 1 {
		panic("invalid cache size")
	}

	return &LFUNoTS{
		frequencyList: list.New(),
		cache:         NewMemoryNoTS(),
		size:          size,
		currentSize:   0,
	}
}

// set a new key-value pair
func (l *LFUNoTS) Get(key string) (interface{}, error) {
	return nil, nil
}

// set a new key-value pair
func (l *LFUNoTS) Set(key string, value interface{}) error {
	return l.set(key, value)
}

func (l *LFUNoTS) set(key string, value interface{}) error {
	res, err := l.cache.Get(key)
	fmt.Println("ERR1")
	if err != nil && err != ErrNotFound {
		return err
	}
	fmt.Println("ERR2")
	if err == ErrNotFound {
		//create new cache item
		fmt.Println("ERR3")
		ci := newCacheItem(key, value)
		l.cache.Set(key, value)
		fmt.Println("BURASI1")
		l.incr(ci)
		fmt.Println("ERR4")
		if l.currentSize > l.size {
			// then evict some data from head of linked list.
			l.evict(l.frequencyList.Front())
		}

	} else {
		//update existing one
		l.cache.Set(key, value)
		l.incr(res.(*cacheItem))
	}

	return nil
}

type entry struct {
	// freqCount holds the frequency number
	freqCount int

	// itemCount holds the items how many exist in list
	listEntry map[*cacheItem]byte
}

// incr increments the usage of cache items
func (l *LFUNoTS) incr(ci *cacheItem) {
	var nextValue int
	var nextPosition *list.Element
	// update existing one
	if ci.freqElement != nil {
		fmt.Println("INC1")
		nextValue = ci.freqElement.Value.(*entry).freqCount + 1
		// replace the position of frequency element
		nextPosition = ci.freqElement.Next()
	} else {
		fmt.Println("INC2")
		// create new frequency element for cache item
		// ci.freqElement is nil so next value of freq will be 1
		nextValue = 1
		// we created new element and its position will be head of linked list
		nextPosition = l.frequencyList.Front()
		l.currentSize++
		fmt.Println("INC2")
	}

	// we need to check position first, otherwise it will panic if we try to fetch value of entry
	if nextPosition == nil || nextPosition.Value.(*entry).freqCount != nextValue {
		fmt.Println("INC3")
		// create new entry node for linked list
		entry := newEntry(nextValue)
		if ci.freqElement == nil {
			nextPosition = l.frequencyList.PushFront(entry)
		} else {
			nextPosition = l.frequencyList.InsertAfter(entry, ci.freqElement)
		}
	}
	nextPosition.Value.(*entry).listEntry[ci] = 1
	ci.freqElement = nextPosition

	if ci.freqElement.Prev() != nil {
		l.remove(ci, ci.freqElement.Prev())
	}
	fmt.Println("LIST IS:", l)
	fmt.Println("FREQUENCY LIST IS:", l.frequencyList)
	fmt.Println("FREQUENCY LIST length IS:", l.frequencyList.Len())
	fmt.Println("CURRENT SIZE IS:", l.currentSize)
}

func (l *LFUNoTS) remove(ci *cacheItem, position *list.Element) {
	entry := position.Value.(*entry).listEntry
	delete(entry, ci)
	if len(entry) == 0 {
		l.frequencyList.Remove(position)
	}
}

func newEntry(freqCount int) *entry {
	return &entry{
		freqCount: freqCount,
		listEntry: make(map[*cacheItem]byte),
	}
}

func newCacheItem(key string, value interface{}) *cacheItem {
	return &cacheItem{
		k: key,
		v: value,
	}

}

func (l *LFUNoTS) Delete(key string) error {
	return nil
}

// evict deletes the element from list
func (l *LFUNoTS) evict(e *list.Element) error {
	if e == nil {
		return nil
	}

	for entry, _ := range e.Value.(*entry).listEntry {
		l.cache.Delete(entry.k)
		l.remove(entry, e)
		l.currentSize--
		break
	}

	return nil
}
