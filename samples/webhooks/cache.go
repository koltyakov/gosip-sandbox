package main

import "sync"

// Cache - in memory cache struct
type Cache struct {
	mx sync.Mutex
	m  map[string]string
}

// NewCache - cache constructor
func NewCache() *Cache {
	return &Cache{
		m: make(map[string]string),
	}
}

// Load loads a value by a key
func (c *Cache) Load(key string) (string, bool) {
	c.mx.Lock()
	defer c.mx.Unlock()
	val, ok := c.m[key]
	return val, ok
}

// Store saves a value for a key
func (c *Cache) Store(key string, value string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.m[key] = value
}
