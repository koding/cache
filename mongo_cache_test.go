package cache

import (
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"
)

var (
	defaultSession = &mgo.Session{}

	// config options for MongoCache
	ttl = func(m *MongoCache) {
		m.TTL = 2 * time.Minute
	}

	collection = func(m *MongoCache) {
		m.CollectionName = "TestCollectionName"
	}

	gcInterval = func(m *MongoCache) {
		m.GCInterval = 2 * time.Minute
	}

	startGC = func(m *MongoCache) {
		m.StartGC = true
	}
)

func TestMongoCacheConfig(t *testing.T) {
	defaultConfig := NewMongoCacheWithTTL(defaultSession)
	if defaultConfig == nil {
		t.Fatal("config should not be nil")
	}
	configTTL := NewMongoCacheWithTTL(defaultSession, ttl)
	if configTTL == nil {
		t.Fatal("ttl config should not be nil")
	}
	if configTTL.TTL != time.Minute*2 {
		t.Fatal("config ttl time should equal 2 minutes")
	}
	config := NewMongoCacheWithTTL(defaultSession, collection, startGC)
	if config == nil {
		t.Fatal("config should not be nil")
	}
	if config.CollectionName != "TestCollectionName" {
		t.Fatal("config collection name should equal 'TestCollectionName'")
	}
	if config.StartGC != true {
		t.Fatal("config StartGC option should not be true")
	}
}
