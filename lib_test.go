package multibtree

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var letterBytes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func genRandomStr() string {

	perm := rand.Perm(len(letterBytes))
	p := ""
	for _, i := range perm[:20] {
		p = p + string(letterBytes[i])
	}
	perm = rand.Perm(len(letterBytes))
	for _, i := range perm[:20] {
		p = p + string(letterBytes[i])
	}
	return p
}

func TestTree(t *testing.T) {

	multiBtree, err := NewMultiTree("./test_files", 20)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer multiBtree.Close()

	now := time.Now()

	for i := 0; i < 20000; i++ {
		err = multiBtree.Insert(genRandomStr())
		if err != nil {
			fmt.Println(err, "write")
			return
		}
	}
	fmt.Println(time.Now().Sub(now))

	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	b, _ := json.Marshal(stats)
	fmt.Println("Last GC was:", string(b))

	now = time.Now()

	var c int = 0

	multiBtree.Scan(func(item string) {
		c++
	})

	runtime.ReadMemStats(stats)
	b, _ = json.Marshal(stats)

	fmt.Println(time.Now().Sub(now))
	fmt.Println("Last GC was:", string(b))
	fmt.Println(c)

}
