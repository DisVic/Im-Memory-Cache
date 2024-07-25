package cache

import (
	"errors"
	"time"
)

type Cache struct {
	data map[string]interface{}
	expireAt time.Time
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
		
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *Cache) Get(key string) (interface{}, error)  {
	value, ok := c.data[key]
	if !ok{
		return nil, errors.New("Не найден")
	}
	return value, nil	 
}

func (c *Cache) Delete(key string) error {
	_, ok := c.data[key]
	if !ok{
		return errors.New("Не найден")
	}
	delete(c.data, key)
	return nil
}

func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration){
	c.data[key] = value
	c.expireAt = time.Now().Add(ttl)
	go c.deleteAfterExpiration(key, ttl)
}

func (c *Cache) deleteAfterExpiration(key string, ttl time.Duration) {
	time.Sleep(ttl)
	_, ok := c.data[key]
	if ok && c.expireAt.Before(time.Now()) {
		delete(c.data, key)
	}
}

