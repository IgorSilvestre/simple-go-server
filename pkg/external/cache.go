package external

import (
	"sync"
	"time"
)

// Cache represents a simple in-memory cache
type Cache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

// cacheItem represents a single item in the cache
type cacheItem struct {
	value      interface{}
	expiration int64
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	cache := &Cache{
		items: make(map[string]cacheItem),
	}
	
	// Start a goroutine to periodically clean up expired items
	go cache.cleanupExpiredItems()
	
	return cache
}

// cleanupExpiredItems removes expired items from the cache
func (c *Cache) cleanupExpiredItems() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for {
		<-ticker.C
		c.mu.Lock()
		now := time.Now().UnixNano()
		for key, item := range c.items {
			if item.expiration > 0 && item.expiration < now {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// Set adds an item to the cache with an optional expiration time
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}
	
	c.items[key] = cacheItem{
		value:      value,
		expiration: expiration,
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	
	// Check if the item has expired
	if item.expiration > 0 && item.expiration < time.Now().UnixNano() {
		return nil, false
	}
	
	return item.value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.items = make(map[string]cacheItem)
}

// GlobalCache is a singleton instance of Cache that can be used throughout the application
var GlobalCache = NewCache()