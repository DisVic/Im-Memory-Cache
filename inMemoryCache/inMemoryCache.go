package cache

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/exp/maps"
)

type Cache[K comparable, T any] struct {
	mu   sync.RWMutex
	data map[K]CacheItem[T]
}

type CacheItem[T any] struct {
	value  T
	expiry time.Time
}

func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		data: make(map[K]CacheItem[T]),
	}
}

func (c *Cache[K, T]) deleteAfterExpiration(key K, ttl time.Duration) {
	time.Sleep(ttl)
	_, ok := c.data[key]
	if ok && c.data[key].expiry.Before(time.Now()) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
	}
}

func (c *Cache[K, T]) Set(key K, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem[T]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
	go c.deleteAfterExpiration(key, ttl)
}

func (c *Cache[K, T]) Get(key K) (T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if value, ok := c.data[key]; !ok {
		return zeroValue[T](), errors.New("не найден")
	} else {
		if value.expiry.Before(time.Now()) {
			delete(c.data, key)
			return zeroValue[T](), errors.New("срок жизни кэша истёк")
		}
		return value.value, nil
	}
}

func zeroValue[T any]() T {
	var value T
	return value
}

func (c *Cache[K, T]) Delete(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.data[key]; !ok {
		return errors.New("не найден")
	}
	delete(c.data, key)
	return nil
}

func (c *Cache[K, T]) Clear() {
	c.mu.Lock()
	maps.Clear(c.data)
	c.mu.Unlock()
}
