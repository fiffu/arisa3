package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type animal string

func (n animal) CacheKey() string { return string(n) }

const (
	cat = animal("meow")
	dog = animal("woof")
)

func Test_memoryCache(t *testing.T) {
	cache := newMemoryCache[animal, string](3 * time.Millisecond)

	// cache miss
	_, ok := cache.Peek("asdasdasdasd")
	assert.False(t, ok)

	// cache put
	cache.Put(dog)
	a, ok := cache.Peek(dog.CacheKey())
	assert.True(t, ok)
	assert.Equal(t, dog, a)

	// cache delete
	cache.Delete(dog.CacheKey())
	_, ok = cache.Peek(dog.CacheKey())
	assert.False(t, ok)

	// cache expiry
	cache.Put(cat)
	cache.clock = func() time.Time { return time.Now().Add(1 * time.Hour) }
	_, ok = cache.Peek(cat.CacheKey())
	assert.False(t, ok)

	// cache drop
	cache.Put(dog)
	cache.Put(cat)
	cache.Drop()
	_, ok = cache.Peek(dog.CacheKey())
	assert.False(t, ok)
	_, ok = cache.Peek(cat.CacheKey())
	assert.False(t, ok)

}
