package cache

import "container/list"

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
	res, err := l.cache.Get(key)
	if err != nil {
		return nil, err
	}

	ci := res.(*cacheItem)

	// increase usage of cache item
	l.incr(ci)
	return ci.v, nil
}

// set a new key-value pair
func (l *LFUNoTS) Set(key string, value interface{}) error {
	return l.set(key, value)
}

func (l *LFUNoTS) set(key string, value interface{}) error {
	res, err := l.cache.Get(key)
	if err != nil && err != ErrNotFound {
		return err
	}

	if err == ErrNotFound {
		//create new cache item
		ci := newCacheItem(key, value)

		// if cache size si reached to max size
		// then first remove lfu item from the list
		if l.currentSize >= l.size {
			// then evict some data from head of linked list.
			l.evict(l.frequencyList.Front())
		}

		l.cache.Set(key, ci)
		l.incr(ci)

	} else {
		//update existing one
		val := res.(*cacheItem)
		val.v = value
		l.cache.Set(key, val)
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
		nextValue = ci.freqElement.Value.(*entry).freqCount + 1
		// replace the position of frequency element
		nextPosition = ci.freqElement.Next()
	} else {
		// create new frequency element for cache item
		// ci.freqElement is nil so next value of freq will be 1
		nextValue = 1
		// we created new element and its position will be head of linked list
		nextPosition = l.frequencyList.Front()
		l.currentSize++
	}

	// we need to check position first, otherwise it will panic if we try to fetch value of entry
	if nextPosition == nil || nextPosition.Value.(*entry).freqCount != nextValue {
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
	res, err := l.cache.Get(key)
	if err != nil && err != ErrNotFound {
		return err
	}

	// we dont need to delete if already doesn't exist
	if err == ErrNotFound {
		return nil
	}

	ci := res.(*cacheItem)

	l.remove(ci, ci.freqElement)
	l.currentSize--
	return l.cache.Delete(key)
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
