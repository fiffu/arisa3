package cardboard

//go:generate mockgen -source=interfaces.go -destination=./interfaces_mock.go -package=cardboard

import (
	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/types"
)

type Actual string
type Alias string
type Filter func(posts []*api.Post) []*api.Post
type TagOperation string

const (
	Promote TagOperation = "promote"
	Demote  TagOperation = "demote"
	Omit    TagOperation = "omit"
	Noop    TagOperation = ""
)

type IDomain interface {
	TagsSearch(s string) ([]*api.TagSuggestion, error)
	PostsSearch(IQueryPosts) ([]*api.Post, error)
	PostsResult(IQueryPosts, []*api.Post) (types.IEmbed, error)

	SetPromote(guildID, tagName string) error
	SetDemote(guildID, tagName string) error
	SetOmit(guildID, tagName string) error
	SetAlias(guildID string, alias Alias, actual Actual) error
	GetAliases(guildID string) (map[Alias]Actual, error)
}

// IQueryPosts is the interface of a query for posts, interpreted within the domain (not the API)
type IQueryPosts interface {
	// Returns whether a query should be executed without pre-processing of tags or
	// post-processing of results.
	MagicMode() bool
	// Single-tag search term
	Term() string
	// Set the Term value
	SetTerm(string)
	// Enumerate tags in the query.
	Tags() []string
	// Optional GuildID that the query originated from.
	GuildID() string
}

type IRepository interface {
	GetAliases(guildID string) (map[Alias]Actual, error)
	SetAlias(guildID string, ali Alias, act Actual) error

	GetTagOperations(guildID string) (map[string]TagOperation, error)
	GetPromotes(guildID string) ([]string, error)
	GetDemotes(guildID string) ([]string, error)
	GetOmits(guildID string) ([]string, error)

	SetPromote(guildID string, tag string) error
	SetDemote(guildID string, tag string) error
	SetOmit(guildID string, tag string) error
}
