package main

import (
	"fmt"
	cache "inMemoryCache/inMemoryCache"
	"time"
)

func main() {
	c := cache.NewCache()
	c.SetWithTTL("kk", "jkew", time.Second * 2)
	fmt.Println(c.Get("kk"))
	time.Sleep(time.Second * 4)
	fmt.Println(c.Get("kk"))
}