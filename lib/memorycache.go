package lib

import "time"

// memorycache.go implements an in-memory cache.

type ICacheable[K comparable] interface {
	CacheKey() K
}

type ICache[T ICacheable[K], K comparable] interface {
	Peek(key K) (T, bool)
	Put(data T)
	Delete(key K)
	Drop()
}

type memoryCache[T ICacheable[K], K comparable] struct {
	data    map[K]T
	dataTTL map[K]time.Time
	expiry  time.Duration
	clock   func() time.Time
}

func NewCache[T ICacheable[K], K comparable](expiry time.Duration) ICache[T, K] {
	return newMemoryCache[T, K](expiry)
}

func newMemoryCache[T ICacheable[K], K comparable](expiry time.Duration) *memoryCache[T, K] {
	return &memoryCache[T, K]{
		data:    make(map[K]T),
		dataTTL: make(map[K]time.Time),
		expiry:  expiry,
		clock:   time.Now,
	}
}

func (c *memoryCache[T, K]) Peek(key K) (t T, ok bool) {
	expiryTime, ok := c.dataTTL[key]
	if !ok {
		return
	}
	if c.clock().After(expiryTime) {
		ok = false
		return
	}

	t, ok = c.data[key]
	return
}

func (c *memoryCache[T, K]) Put(data T) {
	key := data.CacheKey()
	c.dataTTL[key] = time.Now().Add(c.expiry)
	c.data[key] = data
}

func (c *memoryCache[T, K]) Delete(key K) {
	delete(c.dataTTL, key)
	delete(c.data, key)
}

func (c *memoryCache[T, K]) Drop() {
	c.data = make(map[K]T)
	c.dataTTL = make(map[K]time.Time)
}
