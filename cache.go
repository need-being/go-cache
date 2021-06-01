package cache

import (
	"sync"
	"time"
)

type item struct {
	value  interface{}
	expiry time.Time
}

type cache struct {
	store map[string]item
	lock  sync.RWMutex
}

func New() Cache {
	return &cache{
		store: make(map[string]item),
	}
}

func (c *cache) Set(key string, value interface{}, ttl time.Duration) {
	if ttl <= 0 {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.store[key] = item{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.store[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiry) {
		return nil, false
	}
	return item.value, true
}

func (c *cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.store, key)
}
