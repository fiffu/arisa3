package cardboard

import (
	"fmt"
	"time"

	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/lib"
)

const (
	cacheKeyAlias      string = "aliases"
	cacheKeyOperations string = "operations"

	Promote TagOperation = "promote"
	Demote  TagOperation = "demote"
	Omit    TagOperation = "omit"
	Noop    TagOperation = ""
)

type TagOperation string
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

type repo struct {
	db    database.IDatabase
	cache lib.ICache
}

func NewRepository(db database.IDatabase) IRepository {
	return &repo{
		db,
		lib.NewMemoryCache(),
	}
}

func (r *repo) peekAliasesMap() (AliasesMap, bool) {
	if cached, ok := r.cache.Peek(cacheKeyAlias); ok {
		if data, ok := (cached.CacheData()).(AliasesMap); ok {
			return data, ok
		}
	}
	return nil, false
}

func (r *repo) peekOperationsMap() (OperationsMap, bool) {
	if cached, ok := r.cache.Peek(cacheKeyOperations); ok {
		if data, ok := (cached.CacheData()).(OperationsMap); ok {
			return data, ok
		}
	}
	return nil, false
}

func (r *repo) peekTagOperation(oper TagOperation) ([]string, bool) {
	if cached, ok := r.cache.Peek(string(oper)); ok {
		if data, ok := (cached.CacheData()).([]string); ok {
			return data, ok
		}
	}
	return nil, false
}

func (r *repo) GetAliases() (map[Alias]Actual, error) {
	if cached, ok := r.peekAliasesMap(); ok {
		return cached, nil
	}

	aliases := make(map[Alias]Actual)
	rows, err := r.db.Query(
		"SELECT alias, actual FROM aliases",
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ali, act string
		if err := rows.Scan(&ali, &act); err != nil {
			return nil, err
		}
		aliases[Alias(ali)] = Actual(act)
	}

	r.cache.Put(AliasesMap(aliases))

	return aliases, nil
}

func (r *repo) SetAlias(ali Alias, act Actual) error {
	if _, err := r.db.Exec(
		"INSERT INTO aliases(alias, actual) VALUES ($1, $2)",
		string(ali), string(act),
	); err != nil {
		return nil
	}
	r.cache.Delete(cacheKeyAlias)
	return nil
}

func (r *repo) getTagsByOperation(oper TagOperation) ([]string, error) {
	if cached, ok := r.peekTagOperation(oper); ok {
		return cached, nil
	}

	tags := make([]string, 0)
	rows, err := r.db.Query(
		fmt.Sprintf("SELECT tag FROM tag_%s", oper),
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	r.cache.Put(CachedTagList{
		string(oper),
		tags,
	})

	return tags, nil
}

func (r *repo) setTagOperation(tag string, oper TagOperation) error {
	if _, err := r.db.Exec(
		fmt.Sprintf("INSERT INTO tag_%s(tag) VALUES ($1) ON CONFLICT DO NOTHING", string(oper)),
		tag,
	); err != nil {
		return nil
	}
	cacheKey := string(oper)
	r.cache.Delete(cacheKey)
	r.cache.Delete(cacheKeyOperations)
	return nil
}

func (r *repo) GetPromotes() ([]string, error) { return r.getTagsByOperation(Promote) }
func (r *repo) GetDemotes() ([]string, error)  { return r.getTagsByOperation(Demote) }
func (r *repo) GetOmits() ([]string, error)    { return r.getTagsByOperation(Omit) }
func (r *repo) SetPromote(s string) error      { return r.setTagOperation(s, Promote) }
func (r *repo) SetDemote(s string) error       { return r.setTagOperation(s, Demote) }
func (r *repo) SetOmit(s string) error         { return r.setTagOperation(s, Omit) }

func (r *repo) GetTagOperations() (map[string]TagOperation, error) {
	if cached, ok := r.peekOperationsMap(); ok {
		return cached, nil
	}

	mapping := make(map[string]TagOperation)

	for _, oper := range []TagOperation{Promote, Demote, Omit} {
		tags, err := r.getTagsByOperation(oper)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags {
			mapping[tag] = oper
		}
	}

	r.cache.Put(OperationsMap(mapping))
	return mapping, nil
}
