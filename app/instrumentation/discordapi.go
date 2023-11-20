package instrumentation

import (
	"context"
	"regexp"

	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/lib/functional"
)

var (
	discordAPIMatcher = (&apiTreeMatcher{
		Base: `discord.com/api/.+`,
		Resources: []string{
			`/applications/.+/commands`,
			`/guilds/.+/members/.+`,
			`/interactions/.+/.+/callback`,
			`/guilds/.+/roles/.+`,
			`/guilds/.+/roles`,
		},
	}).MustCompile()
)

func MatchDiscordAPIPath(ctx context.Context, path string) string {
	return discordAPIMatcher.Match(ctx, path)
}

type apiTreeMatcher struct {
	Base      string
	Resources []string

	base      *regexp.Regexp
	resources []*regexp.Regexp
}

func (m *apiTreeMatcher) MustCompile() *apiTreeMatcher {
	m.base = regexp.MustCompile(m.Base)
	m.resources = functional.Map(m.Resources, func(resource string) *regexp.Regexp {
		return regexp.MustCompile(m.Base + resource)
	})
	return m
}

func (m *apiTreeMatcher) Match(ctx context.Context, opaquePath string) string {
	base := m.base.FindString(opaquePath)
	if base == "" {
		return ""
	}

	for _, pattern := range m.resources {
		if pattern.MatchString(opaquePath) {
			return pattern.String()
		}
	}
	log.Warnf(ctx, "apiTreeMatcher: no route matched under %s for path: %s", m.Base, opaquePath)
	return opaquePath
}
