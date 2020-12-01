package http

import (
	"github.com/mcache/cache"
	"github.com/mcache/consistendhash"
	"log"
	"net/http"
	"strings"
	"sync"
)

const defaultBasePath = "/cache/"

type Server struct {
	addr     string
	basePath string
	peers    *consistendhash.Map
	sync.Mutex
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		basePath: defaultBasePath,
	}
}

func (s *Server) RegisterPeers(peers ...string) {
	s.Lock()
	defer s.Unlock()

	s.peers = consistendhash.New(50, nil)
	s.peers.AddNodes(peers...)
}

func (s *Server) PickPeer(key string) string {
	s.Lock()
	defer s.Unlock()

	if peer := s.peers.GetNode(key); peer != "" && peer != s.addr {
		return peer
	}
	return ""
}

func (s *Server) Run() {
	log.Printf("cache server is running at %s\n", s.addr)
	log.Fatal(http.ListenAndServe(s.addr, s))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.basePath) {
		http.Error(w, "unknown path", http.StatusBadRequest)
		return
	}
	log.Printf("[Server %s] %s\n", s.addr, r.URL.Path)

	params := strings.SplitN(r.URL.Path[len(s.basePath):], "/", 2)
	if len(params) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	group := cache.GetGroup(params[0])
	if group == nil {
		http.Error(w, "unknown group", http.StatusBadRequest)
		return
	}

	view, err := group.Get(params[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
