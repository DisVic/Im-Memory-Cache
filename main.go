package main

import (
	"fmt"
	cache "inMemoryCache/inMemoryCache"
)

func main() {
	c := cache.NewCache[string, string](10)
	c.Set("age", "12", 100)
	c.Set("name", "wan", 60)
	fmt.Println(c.Get("age"))
	fmt.Println(c.Get("name"))
	c.Delete("name")
	fmt.Println(c.Get("name"))

}
