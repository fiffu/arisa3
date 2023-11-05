package api

import "context"

//go:generate mockgen -source=interfaces.go -destination=./interfaces_mock.go -package=api

type IClient interface {
	UseAuth() bool
	FaviconURL() string
	GetPosts(ctx context.Context, tags []string) ([]*Post, error)
	GetTags(ctx context.Context, tags []string) (map[string]*Tag, error)
	GetTagsMatching(ctx context.Context, pattern string) ([]*Tag, error)

	// Get autocomplete suggestions for a given string.
	AutocompleteTag(ctx context.Context, query string) ([]*TagSuggestion, error)
}
