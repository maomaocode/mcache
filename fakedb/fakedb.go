package fakedb

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

var (
	once   sync.Once
	dbInstance *db
)

type db struct {
	m map[string]string
}

func FakeDB() *db {
	once.Do(func() {
		dbInstance = &db{map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
			"key4": "value4",
			"key5": "value5",
		}}
	})
	return dbInstance
}

func (d *db) Get(key string) (string, bool) {
	v, ok := d.m[key]
	if ok {
		return v, ok
	}
	return "", ok
}

func (d *db) LoaderFunc() func (key string) ([]byte, error) {
	return func(key string) ([]byte, error) {
		db := FakeDB()
		v, ok := db.Get(key)
		log.Printf("[slow db] search key=%s %v\n", key, ok)
		time.Sleep(time.Millisecond * 100 * time.Duration(rand.Intn(10)))
		return []byte(v), nil
	}
}


