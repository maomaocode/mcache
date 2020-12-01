package cache

import (
	"fmt"
	"github.com/mcache/lru"
	"github.com/mcache/peer"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

type Group struct {
	name       string
	loader     LoaderFunc
	mainCache  cache
	peerPicker peer.PeerPicker
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



func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) BindPeerPicker (peerPicker peer.PeerPicker) {
	mu.Lock()
	defer mu.Unlock()
	g.peerPicker = peerPicker
}

func (g *Group) Get(key string) (ByteView, error) {

	if g.peerPicker != nil {
		p := g.peerPicker .PickPeer(key)
		if p != "" {
			log.Printf("load from remote %s\n", p)
			data, err := g.loadFromRemote(p, key)
			if err != nil {
				return ByteView{}, err
			}
			return ByteView{data}, nil
		}
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[cache] hit key=%s", key)
		return v, nil
	}
	return g.load(key)
}

func (g *Group) loadFromRemote(peer, key string) ([]byte, error) {
		u := fmt.Sprintf(
			"%v%v/%v",
			peer,
			g.name,
			key,
		)
		res, err := http.Get(u)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("server returned: %v", res.Status)
		}

		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("reading response body: %v", err)
		}

		return bytes, nil
}

func (g *Group) load(key string, ) (ByteView, error) {
	bytes, err := g.loader.Load(key)
	if err != nil {
		return ByteView{}, err
	}

	view := ByteView{data: cloneBytes(bytes)}
	g.populateCache(key, view)
	return view, nil
}

func (g *Group) populateCache(key string, view ByteView) {
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
