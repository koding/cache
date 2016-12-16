package cache

import "time"

// MongoCache...
type MongoCache struct {
	mongeSession *mgo.Session
	// cache holds the cache data
	cache *KeyValue

	// ttl is a duration for a cache key to expire
	ttl time.Duration
}

func NewMongoCacheWithTTL(session *mgo.Session, ttl time.Duration) *MongoCache {
	return &MongoCache{
		mongeSession: session,
		cache:        &KeyValue{},
		ttl:          ttl,
	}
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
		ExpireAt:  time.Now().UTC().Add(m.ttl),
	}

	return m.CreateKeyValueWithExpiration(kv)
}
