package http

import (
	"testing"
)

func TestConnPool_ServeHTTP(t *testing.T) {
	server := NewServer("localhost:9999", "kv")
	server.Run()
}
