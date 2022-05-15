package cardboard

//go:generate mockgen -source=interfaces.go -destination=./interfaces_mock.go -package=cardboard

import (
	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/types"
)

type Actual string
type Alias string
type Filter func(posts []*api.Post) []*api.Post

type IDomain interface {
	PostsSearch(IQueryPosts) ([]*api.Post, error)
	PostsResult(IQueryPosts, []*api.Post) (types.IEmbed, error)

	PromoteTag(tagName string) error
	DemoteTag(tagName string) error
	OmitTag(tagName string) error
	AliasTag(actual, alias string) error
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
	// Render tags into a search string.
	String() string
}

type IRepository interface {
	GetAliases() (map[Alias]Actual, error)
	SetAlias(Alias, Actual) error

	GetTagOperations() (map[string]TagOperation, error)
	GetPromotes() ([]string, error)
	GetDemotes() ([]string, error)
	GetOmits() ([]string, error)

	SetPromote(string) error
	SetDemote(string) error
	SetOmit(string) error
}
