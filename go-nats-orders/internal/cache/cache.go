package cache

import "sync"


type Cache struct {
	mu   sync.RWMutex
	data map[string]map[string]interface{}
}

func New() *Cache {
	return &Cache{data: make(map[string]map[string]interface{})}
}

func (c *Cache) Set(key string, val map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = val
}

func (c *Cache) Get(key string) (map[string]interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *Cache) GetAll() map[string]map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]map[string]interface{}, len(c.data))
	for k, v := range c.data {
		out[k] = v
	}
	return out
}
