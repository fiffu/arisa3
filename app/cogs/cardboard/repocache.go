package cardboard

import (
	"time"

	"github.com/fiffu/arisa3/lib"
)

type AliasesMap map[Alias]Actual
type OperationsMap map[string]TagOperation

type CachedTagList struct {
	operation string
	list      []string
}

func (c AliasesMap) CacheKey() string             { return cacheKeyAlias }
func (c AliasesMap) CacheData() interface{}       { return c }
func (c AliasesMap) CacheDuration() time.Duration { return 24 * 7 * time.Hour }

func (o OperationsMap) CacheKey() string             { return cacheKeyOperations }
func (o OperationsMap) CacheData() interface{}       { return o }
func (o OperationsMap) CacheDuration() time.Duration { return 24 * 14 * time.Hour } // longer as this is derived cache

func (t CachedTagList) CacheKey() string             { return t.operation }
func (t CachedTagList) CacheData() interface{}       { return t.list }
func (t CachedTagList) CacheDuration() time.Duration { return 24 * 7 * time.Hour }

func (r *repo) newGuildCache(guildID string) lib.ICache {
	cache := lib.NewMemoryCache()
	r.caches[guildID] = cache
	return cache
}

func (r *repo) getGuildCache(guildID string) (lib.ICache, bool) {
	guildCache, ok := r.caches[guildID]
	return guildCache, ok
}

func (r *repo) ensureGuildCache(guildID string) lib.ICache {
	if guildCache, ok := r.getGuildCache(guildID); ok {
		return guildCache
	}
	return r.newGuildCache(guildID)
}

func (r *repo) clearGuildCacheKey(guildID string, cacheKey string) {
	if guildCache, ok := r.getGuildCache(guildID); ok {
		guildCache.Delete(cacheKey)
	}

}

// AliasesMap

func (r *repo) putAliasesMap(guildID string, mapping AliasesMap) {
	r.ensureGuildCache(guildID).Put(mapping)
}

func (r *repo) clearAliasesMap(guildID string) {
	r.clearGuildCacheKey(guildID, (&AliasesMap{}).CacheKey())
}

func (r *repo) peekAliasesMap(guildID string) (AliasesMap, bool) {
	cacheKey := (&AliasesMap{}).CacheKey()
	if guildCache, ok := r.getGuildCache(guildID); ok {
		if cached, ok := guildCache.Peek(cacheKey); ok {
			if data, ok := (cached.CacheData()).(AliasesMap); ok {
				return data, ok
			}
		}
	}
	return nil, false
}

// OperationsMap

func (r *repo) putOperationsMap(guildID string, mapping OperationsMap) {
	r.ensureGuildCache(guildID).Put(mapping)
}

func (r *repo) clearOperationsMap(guildID string) {
	r.clearGuildCacheKey(guildID, (&OperationsMap{}).CacheKey())
}

func (r *repo) peekOperationsMap(guildID string) (OperationsMap, bool) {
	cacheKey := (&OperationsMap{}).CacheKey()
	if guildCache, ok := r.getGuildCache(guildID); ok {
		if cached, ok := guildCache.Peek(cacheKey); ok {
			if data, ok := (cached.CacheData()).(OperationsMap); ok {
				return data, ok
			}
		}
	}
	return nil, false
}

// TagOperation lists

func (r *repo) putTagOperations(guildID string, tagList CachedTagList) {
	r.ensureGuildCache(guildID).Put(tagList)
}

func (r *repo) clearTagOperation(guildID string, oper TagOperation) {
	r.clearGuildCacheKey(guildID, string(oper))
}

func (r *repo) peekTagOperation(guildID string, oper TagOperation) ([]string, bool) {
	cacheKey := string(oper)
	if guildCache, ok := r.getGuildCache(guildID); ok {
		if cached, ok := guildCache.Peek(cacheKey); ok {
			if data, ok := (cached.CacheData()).([]string); ok {
				return data, ok
			}
		}
	}
	return nil, false
}
