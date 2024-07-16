package main

import (
	"fmt"
	cache "inMemoryCache/inMemoryCache"
)

func main() {
	c := cache.NewCache()
	c.Set("age", 12)
	c.Set("name", "wan")
	fmt.Println(c.Get("age"))
	fmt.Println(c.Get("name"))
	c.Delete("name")
	fmt.Println(c.Get("name"))
	
}