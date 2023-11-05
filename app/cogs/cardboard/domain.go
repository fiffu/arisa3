package cardboard

import (
	"context"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/database"
	"github.com/fiffu/arisa3/app/types"
)

type domain struct {
	repo   IRepository
	client api.IClient
}

func NewDomain(db database.IDatabase, cfg *Config) *domain {
	apiClient := api.NewClient(cfg.User, cfg.APIKey, cfg.APITimeoutSecs)
	return &domain{
		NewRepository(db),
		apiClient,
	}
}

func (d *domain) TagsSearch(ctx context.Context, query string) ([]*api.TagSuggestion, error) {
	return d.client.AutocompleteTag(ctx, query)
}

func (d *domain) PostsSearch(ctx context.Context, q IQueryPosts) ([]*api.Post, error) {
	if q.MagicMode() {
		return d.magicSearch(ctx, q, true)
	}
	return d.boringSearch(ctx, q)
}

func (d *domain) PostsResult(ctx context.Context, query IQueryPosts, posts []*api.Post) (types.IEmbed, error) {
	if len(posts) > 0 {
		return d.formatResult(ctx, query, posts)
	} else {
		return d.formatZeroResults(query), nil
	}
}

// Remaining methods proxy to the repo

func (d *domain) SetPromote(ctx context.Context, gid, tagName string) error {
	return d.repo.SetPromote(gid, tagName)
}

func (d *domain) SetDemote(ctx context.Context, gid, tagName string) error {
	return d.repo.SetDemote(gid, tagName)
}

func (d *domain) SetOmit(ctx context.Context, gid, tagName string) error {
	return d.repo.SetOmit(gid, tagName)
}

func (d *domain) SetAlias(ctx context.Context, gid string, al Alias, ac Actual) error {
	return d.repo.SetAlias(gid, al, ac)
}

func (d *domain) GetAliases(ctx context.Context, gid string) (map[Alias]Actual, error) {
	return d.repo.GetAliases(gid)
}
