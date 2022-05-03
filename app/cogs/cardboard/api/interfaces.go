package api

//go:generate mockgen -source=interfaces.go -destination=./interfaces_mock.go -package=api

type IClient interface {
	GetPosts(tags []string) ([]*Post, error)
	GetTags(tags []string) (map[string]*Tag, error)
	GetTagsMatching(pattern string) ([]*Tag, error)
}
