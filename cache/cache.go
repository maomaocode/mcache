package cache

import (
	"github.com/mcache/lru"
	"log"
	"sync"
)

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name      string
	loader    LoaderFunc
	mainCache cache
}

func NewGroup(name string, capacity int64, loader LoaderFunc) *Group {
	if loader == nil {
		panic("nil loader")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		loader:    loader,
		mainCache: cache{cap: capacity},
	}
	groups[name] = g
	return g
}

func GetGroup(name string)  *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) Get(key string) (ByteView, error) {
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[cache] hit key=%s", key)
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string,) (ByteView, error) {
	bytes, err := g.loader.Load(key)
	if err != nil {
		return ByteView{}, err
	}

	view := ByteView{data: cloneBytes(bytes)}
	g.populateCache(key, view)
	return view, nil
}

func (g *Group) populateCache(key string, view ByteView)  {
	g.mainCache.add(key, view)
}

type cache struct {
	sync.RWMutex
	lru *lru.MCache
	cap int64
}

func (c *cache) add(key string, view ByteView) {
	c.Lock()
	defer c.Unlock()

	if c.lru == nil {
		c.lru = lru.NewLruCache(c.cap, nil)
	}

	c.lru.Add(key, view)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.RLock()
	defer c.RUnlock()

	if c == nil || c.lru == nil {
		return ByteView{}, false
	}

	v := c.lru.Get(key)
	if v == nil {
		return ByteView{}, false
	}
	return v.(ByteView), true
}
