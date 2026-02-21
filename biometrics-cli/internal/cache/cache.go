package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Cache struct {
	mu       sync.RWMutex
	items    map[string]CacheItem
	diskPath string
}

type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration time.Time
	CreatedAt  time.Time
}

type CacheConfig struct {
	DiskPath        string
	MaxMemoryMB     int
	TTL             time.Duration
	CleanupInterval time.Duration
}

var (
	globalCache *Cache
	once        sync.Once
)

func New(config *CacheConfig) *Cache {
	once.Do(func() {
		globalCache = &Cache{
			items:    make(map[string]CacheItem),
			diskPath: config.DiskPath,
		}
		if config.DiskPath != "" {
			globalCache.loadFromDisk()
		}
		go globalCache.startCleanup(config.CleanupInterval)
	})
	return globalCache
}

func Get() *Cache {
	if globalCache == nil {
		return New(&CacheConfig{
			DiskPath:        "./cache",
			TTL:             5 * time.Minute,
			CleanupInterval: 1 * time.Minute,
		})
	}
	return globalCache
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(ttl),
		CreatedAt:  time.Now(),
	}

	if c.diskPath != "" {
		c.saveToDisk(key, value, ttl)
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.Expiration) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) GetString(key string) (string, bool) {
	val, ok := c.Get(key)
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

func (c *Cache) GetInt(key string) (int, bool) {
	val, ok := c.Get(key)
	if !ok {
		return 0, false
	}
	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	}
	return 0, false
}

func (c *Cache) GetBool(key string) (bool, bool) {
	val, ok := c.Get(key)
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}

func (c *Cache) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for k := range c.items {
		keys = append(keys, k)
	}
	return keys
}

func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *Cache) startCleanup(interval time.Duration) {
	if interval <= 0 {
		return
	}

	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.cleanup()
		}
	}()
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.items {
		if now.After(item.Expiration) {
			delete(c.items, key)
		}
	}
}

func (c *Cache) saveToDisk(key string, value interface{}, ttl time.Duration) {
	if c.diskPath == "" {
		return
	}

	os.MkdirAll(c.diskPath, 0755)

	data := struct {
		Value     interface{} `json:"value"`
		TTL       int64       `json:"ttl"`
		CreatedAt int64       `json:"created_at"`
	}{
		Value:     value,
		TTL:       int64(ttl),
		CreatedAt: time.Now().Unix(),
	}

	filename := filepath.Join(c.diskPath, fmt.Sprintf("%s.json", key))
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(data)
}

func (c *Cache) loadFromDisk() {
	if c.diskPath == "" {
		return
	}

	os.MkdirAll(c.diskPath, 0755)

	files, err := os.ReadDir(c.diskPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		key := file.Name()[:len(file.Name())-5]
		filename := filepath.Join(c.diskPath, file.Name())

		data := struct {
			Value     interface{} `json:"value"`
			TTL       int64       `json:"ttl"`
			CreatedAt int64       `json:"created_at"`
		}{}

		fileData, err := os.ReadFile(filename)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(fileData, &data); err != nil {
			continue
		}

		expiration := time.Unix(data.CreatedAt, 0).Add(time.Duration(data.TTL) * time.Second)
		if time.Now().Before(expiration) {
			c.items[key] = CacheItem{
				Key:        key,
				Value:      data.Value,
				Expiration: expiration,
				CreatedAt:  time.Unix(data.CreatedAt, 0),
			}
		}
	}
}

type ModelCache struct {
	cache *Cache
	ttl   time.Duration
}

func NewModelCache(ttl time.Duration) *ModelCache {
	return &ModelCache{
		cache: Get(),
		ttl:   ttl,
	}
}

func (m *ModelCache) StoreModelResult(modelID, taskID, result string) {
	key := fmt.Sprintf("model:%s:task:%s", modelID, taskID)
	m.cache.Set(key, result, m.ttl)
}

func (m *ModelCache) GetModelResult(modelID, taskID string) (string, bool) {
	key := fmt.Sprintf("model:%s:task:%s", modelID, taskID)
	return m.cache.GetString(key)
}

func (m *ModelCache) StoreTaskOutput(taskID, output string) {
	key := fmt.Sprintf("task:%s:output", taskID)
	m.cache.Set(key, output, 10*time.Minute)
}

func (m *ModelCache) GetTaskOutput(taskID string) (string, bool) {
	key := fmt.Sprintf("task:%s:output", taskID)
	return m.cache.GetString(key)
}
