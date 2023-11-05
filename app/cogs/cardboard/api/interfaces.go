package api

//go:generate mockgen -source=interfaces.go -destination=./interfaces_mock.go -package=api

type IClient interface {
	UseAuth() bool
	FaviconURL() string
	GetPosts(tags []string) ([]*Post, error)
	GetTags(tags []string) (map[string]*Tag, error)
	GetTagsMatching(pattern string) ([]*Tag, error)

	// Get autocomplete suggestions for a given string.
	AutocompleteTag(query string) ([]*TagSuggestion, error)
}
