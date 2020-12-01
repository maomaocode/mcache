package main

import (
	"flag"
	"fmt"
	"github.com/mcache/cache"
	"github.com/mcache/fakedb"
	hs "github.com/mcache/http"
	"log"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 0, "Geecache server port")
	flag.Parse()


	apiAddr := "http://localhost:9000"
	cacheGroup := cache.NewGroup("score", 2<<10, fakedb.FakeDB().LoaderFunc())

	peers := []string{
		"localhost:9001",
		"localhost:9002",
		"localhost:9003",
	}
	if port == 2 {
		StartAPIServer(apiAddr, cacheGroup)
		StartCacheServer(peers[port], cacheGroup, peers...)
	}else {
		StartCacheServer(peers[port], cacheGroup, peers...)
	}
}

func StartAPIServer(addr string, cacheGroup *cache.Group) {
	http.Handle("/api", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get("key")
		view, err := cacheGroup.Get(key)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Write(view.ByteSlice())
	}))
	log.Println("fontend server is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nil))
}

func StartCacheServer(addr string, cacheGroup *cache.Group, peers ...string) {
	fmt.Println(addr, peers)
	server := hs.NewServer(addr)
	server.RegisterPeers(peers...)
	cacheGroup.BindPeerPicker(server)
	server.Run()
}
