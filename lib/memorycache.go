package lib

// memorycache.go implements an in-memory cache.

import "time"

type ICache interface {
	Peek(key string) (ICacheable, bool)
	Put(data ICacheable)
	Delete(key string)
}

type ICacheable interface {
	CacheKey() string
	CacheData() interface{}
	CacheDuration() time.Duration
}

// memoryCache implements ICache
type memoryCache map[string]*cacheRecord
type cacheRecord struct {
	data   ICacheable
	expiry time.Time
}

func NewMemoryCache() ICache {
	cache := memoryCache(make(map[string]*cacheRecord))
	return &cache
}

func (c *memoryCache) Peek(key string) (ICacheable, bool) {
	record, ok := (*c)[key]
	if !ok {
		return nil, false
	}
	if time.Now().After(record.expiry) {
		// cache expired
		delete(*c, key)
		return nil, false
	}
	return record.data, true
}

func (c *memoryCache) Put(data ICacheable) {
	key := data.CacheKey()
	duration := data.CacheDuration()
	expiry := time.Now().Add(duration)
	(*c)[key] = &cacheRecord{data, expiry}
}

func (c *memoryCache) Delete(key string) {
	delete(*c, key)
}
