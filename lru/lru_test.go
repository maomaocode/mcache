package lru

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type fakeValue struct {
	data string
}

func (v fakeValue) Len() int {
	return len(v.data)
}

func TestMCache(t *testing.T) {
	lruCache := NewLruCache(21, nil)

	lruCache.Add("key1", fakeValue{"123456"}) // len = 10
	lruCache.Add("key2", fakeValue{"234567"})

	if v := lruCache.Get("key1"); v == nil || v.(fakeValue).data != "123456" {
		t.Fatalf("cache hit fail")
	}
	lruCache.Add("key3", fakeValue{"345678"})

	if v := lruCache.Get("key2"); v != nil {
		t.Fatalf("cache miss fail")
	}
	if v := lruCache.Get("key3"); v == nil || v.(fakeValue).data != "345678" {
		t.Fatalf("cache hit fail")
	}
}

func TestMCacheConcurrent(t *testing.T) {

	lruCache := NewLruCache(21, nil)

	for i := 0; i < 5000; i++ {
		go func(n int) {
			lruCache.Add(strconv.Itoa(rand.Intn(10)), fakeValue{strconv.Itoa(n)})
			fmt.Println("add ", n)
		}(i)
		go func(n int) {
			lruCache.Get(strconv.Itoa(rand.Intn(10)))
			fmt.Println("get ", n)
		}(i)
	}
	time.Sleep(time.Second * 100)
}
