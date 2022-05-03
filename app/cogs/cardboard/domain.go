package cardboard

import (
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

func (d *domain) PostsSearch(q IQueryPosts) ([]*api.Post, error) {
	if q.MagicMode() {
		return d.magicSearch(q, true)
	}
	return d.boringSearch(q)
}

func (d *domain) PostsResult(query IQueryPosts, posts []*api.Post) (types.IEmbed, error) {
	if len(posts) > 0 {
		return d.formatResult(query, posts)
	} else {
		return d.formatZeroResults(query), nil
	}
}
