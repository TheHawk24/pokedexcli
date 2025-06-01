package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entry map[string]cacheEntry
	mu    *sync.Mutex
}

func (c *Cache) Add(key string, value []byte) {
	cache_entry := cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	c.mu.Lock()
	c.entry[key] = cache_entry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	value, ok := c.entry[key]
	c.mu.Unlock()
	if !ok {
		return nil, false
	}
	return value.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	for {
		ticker := time.NewTicker(interval)
		t := <-ticker.C
		for key, value := range c.entry {
			if value_time := value.createdAt; t.Sub(value_time) > interval {
				c.mu.Lock()
				delete(c.entry, key)
				c.mu.Unlock()
			}
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		entry: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return &cache
}
