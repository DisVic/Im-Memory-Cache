package cache

import (
	"errors"
	"sync"
	"time"

	"golang.org/x/exp/maps"
)

type Cache[K comparable, T any] struct {
	mu       sync.RWMutex
	data     map[K]CacheItem[T]
	interval time.Duration
	stop     chan struct{}
}

type CacheItem[T any] struct {
	value  T
	expiry time.Time
}

func NewCache[K comparable, T any](cleanupInterval time.Duration) *Cache[K, T] {
	cache := &Cache[K, T]{
		data:     make(map[K]CacheItem[T]),
		interval: cleanupInterval,
		stop:     make(chan struct{}),
	}
	go cache.cleanupExpiredItems()
	return cache
}

func (c *Cache[K, T]) cleanupExpiredItems() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.mu.Lock()
			for key, item := range c.data {
				if item.expiry.Before(now) {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		case <-c.stop:
			return
		}
	}
}

func (c *Cache[K, T]) Set(key K, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = CacheItem[T]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
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
