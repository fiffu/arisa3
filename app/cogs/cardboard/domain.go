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

func (d *domain) PromoteTag(tagName string) error {
	return d.repo.SetPromote(tagName)
}

func (d *domain) DemoteTag(tagName string) error {
	return d.repo.SetDemote(tagName)
}

func (d *domain) OmitTag(tagName string) error {
	return d.repo.SetOmit(tagName)
}

func (d *domain) AliasTag(actual, alias string) error {
	return d.repo.SetAlias(Alias(alias), Actual(actual))
}
