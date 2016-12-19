package cache

import (
	"testing"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	defaultSession = &mgo.Session{}

	session = initMongo()

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

func initMongo() *mgo.Session {
	ses, err := mgo.Dial("127.0.0.1:27017/test")
	if err != nil {
		panic(err)
	}

	ses.SetSafe(&mgo.Safe{})
	ses.SetMode(mgo.Strong, true)

	return ses
}

func TestMongoCacheConfig(t *testing.T) {
	defaultConfig := NewMongoCacheWithTTL(session)
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

func TestMongoCacheSetOptionFuncs(t *testing.T) {
	defaultConfig := NewMongoCacheWithTTL(defaultSession)
	if defaultConfig == nil {
		t.Fatal("config should not be nil")
	}

	duration := time.Minute * 3
	configTTL := NewMongoCacheWithTTL(defaultSession, SetTTL(duration))
	if configTTL == nil {
		t.Fatal("ttl config should not be nil")
	}
	if configTTL.TTL != duration {
		t.Fatal("config ttl time should equal 2 minutes")
	}

	// check multiple options
	collName := "testingCollectionName"
	config := NewMongoCacheWithTTL(defaultSession, SetCollectionName(collName), SetGCInterval(duration), EnableStartGC())
	if config == nil {
		t.Fatal("config should not be nil")
	}
	if config.CollectionName != collName {
		t.Fatal("config collection name should equal 'TestCollectionName'")
	}
	if config.StartGC != true {
		t.Fatal("config StartGC option should not be true")
	}

	if config.GCInterval != duration {
		t.Fatal("config GCInterval option should equal", duration)
	}
}

func TestMongoCacheGet(t *testing.T) {
	defaultConfig := NewMongoCacheWithTTL(session)
	if defaultConfig == nil {
		t.Fatal("config should not be nil")
	}

	_, err := defaultConfig.Get("test")
	if err != mgo.ErrNotFound {
		t.Fatal("error is:", err)
	}
}

func TestMongoCacheSet(t *testing.T) {
	defaultConfig := NewMongoCacheWithTTL(session)
	if defaultConfig == nil {
		t.Fatal("config should not be nil")
	}
	key := bson.NewObjectId().Hex()
	value := bson.NewObjectId().Hex()

	err := defaultConfig.Set(key, value)
	if err != nil {
		t.Fatal("error should be nil:", err)
	}
	data, err := defaultConfig.Get(key)
	if err != nil {
		t.Fatal("error should be nil:", err)
	}
	if data == nil {
		t.Fatal("data should not be nil")
	}
	if data != value {
		t.Fatal("data should equal:", value, ", but got:", data)
	}
}

func TestMongoCacheDelete(t *testing.T) {
	defaultConfig := NewMongoCacheWithTTL(session)
	if defaultConfig == nil {
		t.Fatal("config should not be nil")
	}
	key := bson.NewObjectId().Hex()
	value := bson.NewObjectId().Hex()

	err := defaultConfig.Set(key, value)
	if err != nil {
		t.Fatal("error should be nil:", err)
	}
	data, err := defaultConfig.Get(key)
	if err != nil {
		t.Fatal("error should be nil:", err)
	}
	if data == nil {
		t.Fatal("data should not be nil")
	}
	if data != value {
		t.Fatal("data should equal:", value, ", but got:", data)
	}
	err = defaultConfig.Delete(key)
	if err != nil {
		t.Fatal("err should be nil, but got", err)
	}
}

func TestMongoCacheTTL(t *testing.T) {
	// duration specifies the time duration to hold the data in mongo
	// after the duration interval, data will be deleted from mongoDB
	duration := time.Second * 20
	defaultConfig := NewMongoCacheWithTTL(session, SetTTL(duration))
	if defaultConfig == nil {
		t.Fatal("config should not be nil")
	}
	key := bson.NewObjectId().Hex()
	value := bson.NewObjectId().Hex()

	err := defaultConfig.Set(key, value)
	if err != nil {
		t.Fatal("error should be nil:", err)
	}
	data, err := defaultConfig.Get(key)
	if err != nil {
		t.Fatal("error should be nil:", err)
	}
	if data != value {
		t.Fatal("data should equal:", value, ", but got:", data)
	}
	time.Sleep(duration)
	_, err = defaultConfig.Get(key)
	if err != mgo.ErrNotFound {
		t.Fatal("error should equal", mgo.ErrNotFound, " but got:", err)
	}
}
