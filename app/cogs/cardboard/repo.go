package cardboard

import (
	"fmt"

	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/lib"
)

const (
	cacheKeyAlias      string = "aliases"
	cacheKeyOperations string = "operations"
)

type repo struct {
	db     database.IDatabase
	caches map[string]lib.ICache // one cache instance per guildID
}

func NewRepository(db database.IDatabase) IRepository {
	return &repo{
		db,
		make(map[string]lib.ICache),
	}
}

func (r *repo) GetAliases(guildID string) (map[Alias]Actual, error) {
	if cached, ok := r.peekAliasesMap(guildID); ok {
		return cached, nil
	}

	aliases := make(map[Alias]Actual)
	rows, err := r.db.Query(
		"SELECT alias, actual FROM aliases WHERE guildid=$1",
		guildID,
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

	r.putAliasesMap(guildID, AliasesMap(aliases))
	return aliases, nil
}

func (r *repo) SetAlias(guildID string, ali Alias, act Actual) error {
	if _, err := r.db.Exec(
		"INSERT INTO aliases(alias, actual, guildid) VALUES ($1, $2, $3)",
		string(ali), string(act), guildID,
	); err != nil {
		return err
	}
	r.clearAliasesMap(guildID)
	return nil
}

func (r *repo) getTagsByOperation(guildID string, oper TagOperation) ([]string, error) {
	if cached, ok := r.peekTagOperation(guildID, oper); ok {
		return cached, nil
	}

	tags := make([]string, 0)
	rows, err := r.db.Query(
		fmt.Sprintf("SELECT tag FROM tag_%s WHERE guildid=$1", oper),
		guildID,
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

	r.putTagOperations(
		guildID,
		CachedTagList{string(oper), tags},
	)

	return tags, nil
}

func (r *repo) setTagOperation(guildID string, tag string, oper TagOperation) error {
	if _, err := r.db.Exec(
		fmt.Sprintf(`INSERT INTO tag_%s(tag, guildid) VALUES ($1, $2) ON CONFLICT DO NOTHING`, string(oper)),
		tag, guildID,
	); err != nil {
		return err
	}
	r.clearTagOperation(guildID, oper)
	r.clearOperationsMap(guildID)
	return nil
}

func (r *repo) GetPromotes(gid string) ([]string, error) { return r.getTagsByOperation(gid, Promote) }
func (r *repo) GetDemotes(gid string) ([]string, error)  { return r.getTagsByOperation(gid, Demote) }
func (r *repo) GetOmits(gid string) ([]string, error)    { return r.getTagsByOperation(gid, Omit) }
func (r *repo) SetPromote(gid, s string) error           { return r.setTagOperation(gid, s, Promote) }
func (r *repo) SetDemote(gid, s string) error            { return r.setTagOperation(gid, s, Demote) }
func (r *repo) SetOmit(gid, s string) error              { return r.setTagOperation(gid, s, Omit) }

func (r *repo) GetTagOperations(guildID string) (map[string]TagOperation, error) {
	if cached, ok := r.peekOperationsMap(guildID); ok {
		return cached, nil
	}

	mapping := make(map[string]TagOperation)

	for _, oper := range []TagOperation{Promote, Demote, Omit} {
		tags, err := r.getTagsByOperation(guildID, oper)
		if err != nil {
			return nil, err
		}
		for _, tag := range tags {
			mapping[tag] = oper
		}
	}

	r.putOperationsMap(guildID, OperationsMap(mapping))
	return mapping, nil
}
