package cache

import (
	"errors"

	"golang.org/x/exp/maps"
)

type Cache struct {
	data map[string]interface{}
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, error) {
	value, ok := c.data[key]
	if !ok {
		return nil, errors.New("не найден")
	}
	return value, nil
}

func (c *Cache) Delete(key string) error {
	_, ok := c.data[key]
	if !ok {
		return errors.New("не найден")
	}
	delete(c.data, key)
	return nil
}

func (c *Cache) Clear() {
	maps.Clear(c.data)
}
