package lib

// memorycache.go implements an in-memory cache.

import "time"

type ICache interface {
	Peek(key string) ICacheable
	Put(data ICacheable)
}

type ICacheable interface {
	CacheKey() string
	CacheDuration() time.Duration
	CacheData() interface{}
}

// memoryCache implements ICache
type memoryCache map[string]*cacheRecord
type cacheRecord struct {
	data   ICacheable
	expiry time.Time
}

func NewMemoryCache() ICache {
	return new(memoryCache)
}

func (c *memoryCache) Peek(key string) ICacheable {
	record, ok := (*c)[key]
	if !ok {
		return nil
	}
	if time.Now().After(record.expiry) {
		// cache expired
		delete(*c, key)
		return nil
	}
	return record.data
}

func (c *memoryCache) Put(data ICacheable) {
	key := data.CacheKey()
	duration := data.CacheDuration()
	expiry := time.Now().Add(duration)
	(*c)[key] = &cacheRecord{data, expiry}
}
