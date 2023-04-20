package cardboard

import (
	"time"

	"github.com/fiffu/arisa3/lib"
)

type guildKey string

const (
	aliasKey      guildKey = "aliases"
	operationsKey guildKey = "operations"
)

// AliasesMap indicates all aliases defined by a guild.
// Cache hierarchy: global -> guildID -> AliasMap
type AliasesMap map[Alias]Actual

// OperationsMap indicates all tags defined by a guild to receive any kind of operation.
// This mapping is derived from a union of all TagsPerOperation.
// Cache hierarchy: global -> guildID -> OperationsMap
type OperationsMap map[string]TagOperation

// TagsPerOperation indicates the tags defined by a guild to receive a particular operation.
// Cache hierarchy: global -> guildID -> operation -> TagsPerOperation
type TagsPerOperation struct {
	op   TagOperation
	list []string
}

func (c AliasesMap) CacheKey() guildKey           { return aliasKey }
func (t TagsPerOperation) CacheKey() TagOperation { return t.op }
func (o OperationsMap) CacheKey() guildKey        { return operationsKey }

func (r *repo) newGuildCache(guildID string) perGuildCache {
	cache := perGuildCache{
		aliases:    lib.NewCache[AliasesMap, guildKey](7 * 24 * time.Hour),
		ops2tags:   lib.NewCache[TagsPerOperation, TagOperation](7 * 24 * time.Hour),
		operations: lib.NewCache[OperationsMap, guildKey](14 * 24 * time.Hour), // derived from ops2tags
	}
	r.caches[guildID] = cache
	return cache
}

func (r *repo) getGuildCache(guildID string) (perGuildCache, bool) {
	guildCache, ok := r.caches[guildID]
	return guildCache, ok
}

func (r *repo) ensureGuildCache(guildID string) perGuildCache {
	if guildCache, ok := r.caches[guildID]; ok {
		return guildCache
	}
	return r.newGuildCache(guildID)
}

// AliasesMap

func (r *repo) putAliasesMap(guildID string, mapping AliasesMap) {
	r.ensureGuildCache(guildID).aliases.Put(mapping)
}

func (r *repo) clearAliasesMap(guildID string) {
	if guildCache, ok := r.caches[guildID]; ok {
		guildCache.aliases.Drop()
	}
}

func (r *repo) peekAliasesMap(guildID string) (AliasesMap, bool) {
	if guildCache, ok := r.caches[guildID]; ok {
		return guildCache.aliases.Peek(aliasKey)
	}
	return nil, false
}

// OperationsMap

func (r *repo) putOperationsMap(guildID string, mapping OperationsMap) {
	r.ensureGuildCache(guildID).operations.Put(mapping)
}

func (r *repo) clearOperationsMap(guildID string) {
	if guildCache, ok := r.caches[guildID]; ok {
		guildCache.operations.Drop()
	}
}

func (r *repo) peekOperationsMap(guildID string) (OperationsMap, bool) {
	if guildCache, ok := r.caches[guildID]; ok {
		return guildCache.operations.Peek(aliasKey)
	}
	return nil, false
}

// TagsPerOperation

func (r *repo) putTagOperations(guildID string, tagList TagsPerOperation) {
	r.ensureGuildCache(guildID).ops2tags.Put(tagList)
}

func (r *repo) clearTagOperation(guildID string, oper TagOperation) {
	if guildCache, ok := r.caches[guildID]; ok {
		guildCache.ops2tags.Drop()
	}
}

func (r *repo) peekTagOperation(guildID string, oper TagOperation) ([]string, bool) {
	if guildCache, ok := r.getGuildCache(guildID); ok {
		if tags, ok := guildCache.ops2tags.Peek(oper); ok {
			return tags.list, true
		}
	}
	return nil, false
}
