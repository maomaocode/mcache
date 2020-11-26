package consistendhash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	m := New(2, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	m.AddNodes("1", "3", "5")
	// 1, 3, 5, 11, 13, 15

	fmt.Println(m.GetNode("6")) // should be 1
	m.AddNodes("7") // add 7, 17
	fmt.Println(m.GetNode("6")) // should be 7

}
