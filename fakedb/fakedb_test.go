package fakedb

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestFakeDB(t *testing.T) {
	db := FakeDB()

	for i := 0; i < 10000; i++ {
		go func() {
			key := fmt.Sprintf("key%d", rand.Intn(6))
			fmt.Println(db.Get(key))
		}()

	}



	time.Sleep(time.Second * 100)

}
