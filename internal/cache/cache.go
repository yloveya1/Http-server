package cache

import (
	"errors"
	"sync"
	"time"
)

type CacheItem struct {
	value        interface{}
	timeCreation time.Time
	timeDuration time.Duration
}

type Cache struct {
	cache map[string]CacheItem
	mu    *sync.Mutex
}

func New() *Cache {
	cacheNew := Cache{
		cache: make(map[string]CacheItem),
		mu:    new(sync.Mutex),
	}
	go cacheNew.TimeExpireTask()
	return &cacheNew
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = CacheItem{
		value:        value,
		timeCreation: time.Now(),
		timeDuration: ttl,
	}
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k, ok := c.cache[key]
	if !ok {
		return nil, errors.New("invalid key")
	}

	return k.value, nil
}

func (c *Cache) Delete(key string) {

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *Cache) TimeExpireTask() {
	for {
		c.DeleteExpired()
		time.Sleep(time.Second)
	}
}

func (c *Cache) DeleteExpired() {
	for key, value := range c.cache {
		if time.Since(value.timeCreation) > value.timeDuration {
			c.Delete(key)
		}
	}
}
