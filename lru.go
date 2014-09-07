package cache

import (
	"container/list"
)

type LRUCache struct {
	list  *list.List
	items Cache
	size  int
}

type kv struct {
	k string
	v interface{}
}

func NewLRUCache(size int) *LRUCache {
	if size < 1 {
		panic("invalid cache size")
	}

	return &LRUCache{
		list:  list.New(),
		items: NewMemory(),
		size:  size,
	}
}

func (l *LRUCache) Get(key string) (interface{}, error) {
	res, err := l.items.Get(key)
	if err != nil {
		return nil, err
	}

	elem := res.(*list.Element)
	// move found item to the head
	l.list.MoveToFront(elem)

	return elem.Value.(*kv).v, nil
}

func (l *LRUCache) Set(key string, val interface{}) error {
	res, err := l.items.Get(key)
	if err != nil && err != ErrNotFound {
		return err
	}

	var elem *list.Element

	// if elem is not in the cache, set it
	if err == ErrNotFound {
		elem = l.list.PushFront(&kv{k: key, v: val})
	} else {
		elem = res.(*list.Element)

		// update the  data
		elem.Value.(*kv).v = val

		// item already exists, so move it to the front of the list
		l.list.MoveToFront(elem)
	}

	err = l.items.Set(key, elem)
	if err != nil {
		return err
	}

	// if the cache is full, evict last LRU entry
	if l.list.Len() > l.size {
		// remove last element from cache
		return l.removeElem(l.list.Back())
	}

	return nil
}

func (l *LRUCache) Delete(key string) error {
	res, err := l.items.Get(key)
	if err != nil && err != ErrNotFound {
		return err
	}

	// item already deleted
	if err == ErrNotFound {
		// surpress not found errors
		return nil
	}

	elem := res.(*list.Element)

	return l.removeElem(elem)
}

func (l *LRUCache) removeElem(e *list.Element) error {
	l.list.Remove(e)
	return l.items.Delete(e.Value.(*kv).k)
}
