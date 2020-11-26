package http

import (
	"github.com/mcache/cache"
	"github.com/mcache/fakedb"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/cache/"

type Server struct {
	addr      string
	basePath  string
	groupName string
}

func NewServer(addr string, groupName string) *Server {
	return &Server{
		addr:      addr,
		basePath:  defaultBasePath,
		groupName: groupName,
	}
}

func(s *Server) Run() {
	cache.NewGroup(s.groupName, 2<<10, fakedb.FakeDB().LoaderFunc())
	log.Fatal(http.ListenAndServe(s.addr, s))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, s.basePath) {
		http.Error(w, "unknown path", http.StatusBadRequest)
		return
	}
	log.Printf("[Server %s] %s\n", r.Method, r.URL.Path)

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
