package lru

import (
	"container/list"
	"log"
)

type Value interface {
	Len() int
}

// MCache is a lru cache, it is not safe for concurrent access
type MCache struct {
	maxBytes  int64      // the max bytes memory can save
	usedBytes int64      // the used memory
	ll        *list.List // double linked-list

	cache map[string]*list.Element // indexes for all elements of ll

	// optional and executed when an entry is purged
	OnEvicted func(key string, value Value)
}

func NewLruCache(maxMemory int64, onEvicted func(key string, value Value)) *MCache {
	if onEvicted == nil {
		onEvicted = func(key string, value Value) {
			log.Printf("[mcache] Evict key:%q value:%q", key, value)
		}
	}

	return &MCache{
		maxBytes:  maxMemory,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

type mEntry struct {
	key   string
	value Value
}

func (c *MCache) Get(key string) Value {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		entry := ele.Value.(*mEntry)
		return entry.value
	}
	return nil
}

func (c *MCache) Add(key string, value Value) {
	if int64(len(key)+value.Len()) > c.maxBytes {
		log.Printf("[mcache] too large to save key: %q value: %q, len: %d\n", key, value, value.Len())
		return
	}

	if ele, ok := c.cache[key]; ok { // update
		c.ll.MoveToFront(ele)
		entry := ele.Value.(*mEntry)
		c.usedBytes += int64(value.Len() - entry.value.Len())
		entry.value = value
	} else {
		ele := c.ll.PushFront(&mEntry{key, value})
		c.usedBytes += int64(len(key) + value.Len())
		c.cache[key] = ele
	}

	for c.usedBytes > 0 && c.maxBytes < c.usedBytes {
		c.removeOldest()
	}
}

func (c *MCache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		entry := ele.Value.(*mEntry)
		delete(c.cache, entry.key)
		c.usedBytes -= int64(len(entry.key) + entry.value.Len())
		c.OnEvicted(entry.key, entry.value)
	}
}

func (c *MCache) Len() int {
	return c.ll.Len()
}
