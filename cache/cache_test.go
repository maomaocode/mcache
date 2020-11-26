package cache

import (
	"fmt"
	"github.com/mcache/fakedb"
	"testing"
)

func TestGroup_Get(t *testing.T) {
	g := NewGroup("group1", 2<<10, fakedb.FakeDB().LoaderFunc())

	v, _ := g.Get("key1")
	fmt.Println(v)

	v, _ = g.Get("key1")
	fmt.Println(v)

	v, _ = g.Get("key2")
	fmt.Println(v)

	v, _ = g.Get("key3")
	fmt.Println(v)

	v, _ = g.Get("key2")
	fmt.Println(v)

	v, _ = g.Get("key3")
	fmt.Println(v)

	v, _ = g.Get("key7")
	fmt.Println(v)
}
