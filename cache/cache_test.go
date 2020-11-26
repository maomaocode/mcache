package cache

import (
	"fmt"
	"github.com/mcache/fakedb"
	"math/rand"
	"testing"
	"time"
)

func TestGroup_Get(t *testing.T) {
	g := NewGroup("group1", 2<<10, fakedb.FakeDB().LoaderFunc())

	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < 50; i++ {
				_, _ = g.Get(fmt.Sprintf("key%d", rand.Intn(6)))
				time.Sleep(time.Millisecond * 500)
			}
		}()
	}

	time.Sleep(time.Second * 20)
}
