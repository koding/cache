package cache

import (
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
)

// MongoCache holds the cache values that will be stored in mongoDB
type MongoCache struct {
	// mongeSession specifies the mongoDB connection
	mongeSession *mgo.Session

	// CollectionName speficies the optional collection name for mongoDB
	// if CollectionName is not set, then default value will be set
	CollectionName string

	// ttl is a duration for a cache key to expire
	TTL time.Duration

	// GCInterval specifies the time duration for garbage collector time interval
	GCInterval time.Duration

	// StartGC starts the garbage collector and deletes the
	// expired keys from mongo with given time interval
	StartGC bool

	// gcTicker controls gc intervals
	gcTicker *time.Ticker

	// done controls sweeping goroutine lifetime
	done chan struct{}

	// Mutex is used for handling the concurrent
	// read/write requests for cache
	sync.RWMutex
}

// Option sets the options specified.
type Option func(*MongoCache)

// NewMongoCacheWithTTL creates a caching layer backed by mongo. TTL's are
// managed either by a background cleaner or document is removed on the Get
// operation. Mongo TTL indexes are not utilized since there can be multiple
// systems using the same collection with different TTL values.
//
// The responsibility of stopping the GC process belongs to the user.
//
// Session is not closed while stopping the GC.
// This function satisfy you to not pass nil value to the function as parameter
// e.g (usage) :
// configure with defaults, just call;
// NewMongoCacheWithTTL(session)
//
// configure ttl duration with;
// NewMongoCacheWithTTL(session, func(m *MongoCache) {
// m.TTL = 2 * time.Minute
// })
//
// configure collection name with;
// NewMongoCacheWithTTL(session, func(m *MongoCache) {
// m.CollectionName = "MongoCacheCollectionName"
// })
func NewMongoCacheWithTTL(session *mgo.Session, configs ...Option) *MongoCache {
	mc := &MongoCache{
		mongeSession:   session,
		TTL:            defaultExpireDuration,
		CollectionName: defaultKeyValueColl,
		GCInterval:     time.Minute,
		StartGC:        false,
	}

	for _, configFunc := range configs {
		configFunc(mc)
	}

	if mc.StartGC {
		mc.StartGCollector(mc.GCInterval)
	}

	return mc
}

// EnableStartGC enables the garbage collector in MongoCache struct
func EnableStartGC() Option {
	return optionStartGC(true)
}

// DisableStartGC disables the garbage collector in MongoCache struct
func DisableStartGC() Option {
	return optionStartGC(false)
}

// optionStartGC chooses the garbage collector option in MongoCache struct
func optionStartGC(b bool) Option {
	return func(m *MongoCache) {
		m.StartGC = b

	}
}

// SetTTL sets the ttl duration in MongoCache as option
func SetTTL(duration time.Duration) Option {
	return func(m *MongoCache) {
		m.TTL = duration

	}
}

// SetGCInterval sets the garbage collector interval in MongoCache struct as option
func SetGCInterval(duration time.Duration) Option {
	return func(m *MongoCache) {
		m.GCInterval = duration

	}
}

// SetCollectionName sets the collection name for mongoDB in MongoCache struct as option
func SetCollectionName(collName string) Option {
	return func(m *MongoCache) {
		m.CollectionName = collName

	}
}

// WithStartGC adds the given value to the WithStartGC
// this is an external way to change WithStartGC value as true
// StartGC option is false as default
// usage:
// NewMongoCacheWithTTL(config).WithStartGC(true)
//
// recommended way is :
// to enable StartGC, use EnableStartGC
// add EnableStartGC as option NewMongoCacheWithTTL(&mgo.session{}, EnableStartGC())
func (m *MongoCache) WithStartGC(isStart bool) *MongoCache {
	m.StartGC = isStart
	return m
}

// Get returns a value of a given key if it exists
func (m *MongoCache) Get(key string) (interface{}, error) {
	return m.GetKeyWithExpireCheck(key)
}

// Set will persist a value to the cache or
// override existing one with the new one
func (m *MongoCache) Set(key string, value interface{}) error {
	return m.set(key, value)
}

// Delete deletes a given key if exists
func (m *MongoCache) Delete(key string) error {
	return m.DeleteKey(key)
}

func (m *MongoCache) set(key string, value interface{}) error {
	kv := &KeyValue{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		ExpireAt:  time.Now().UTC().Add(m.TTL),
	}

	return m.CreateKeyValueWithExpiration(kv)
}

// StartGCollector starts the garbage collector with given time interval
// The expired data will be checked & deleted with given interval time
func (m *MongoCache) StartGCollector(gcInterval time.Duration) {
	if gcInterval <= 0 {
		return
	}

	ticker := time.NewTicker(gcInterval)
	done := make(chan struct{})

	m.Lock()
	m.gcTicker = ticker
	m.done = done
	m.Unlock()

	go func() {
		for {
			select {
			case <-ticker.C:
				m.Lock()
				m.DeleteExpiredKeys()
				m.Unlock()
			case <-done:
				return
			}
		}
	}()
}

// StopGC stops sweeping goroutine.
func (r *MemoryTTL) StopGCol() {
	if r.gcTicker != nil {
		r.Lock()
		r.gcTicker.Stop()
		r.gcTicker = nil
		close(r.done)
		r.done = nil
		r.Unlock()
	}
}
