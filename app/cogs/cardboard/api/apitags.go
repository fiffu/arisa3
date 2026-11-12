package api

import (
	"context"
	"net/url"

	"github.com/fiffu/arisa3/app/log"
)

type Tag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"post_count"`
}

type TagSuggestion struct {
	Name       string // actual, normalized name of tag
	Antecedent string // populated if the suggestion matched on this tag's alias
	PostCount  int
	Link       string // URL to posts with this tag
}

type autocompleteResult struct {
	Value      string `json:"value"`
	Antecedent string `json:"antecedent"`
	PostCount  int    `json:"post_count"`
}

// AutocompleteTag implements api.IClient.
// It uses the JSON form of the undocumented /autocomplete endpoint, because Cloudflare blocks the HTML form for the bot.
func (c *client) AutocompleteTag(ctx context.Context, query string) ([]*TagSuggestion, error) {
	ctx, cancel := c.httpContext(ctx)
	defer cancel()

	var rows []autocompleteResult
	builder := c.autocompleteEndpoint().
		Param("search[query]", query).
		Param("search[type]", "tag_query").
		Param("limit", "10").
		ToJSON(&rows)
	if err := c.fetch(ctx, builder); err != nil {
		return nil, err
	}

	result := make([]*TagSuggestion, 0, len(rows))
	for _, row := range rows {
		result = append(result, &TagSuggestion{
			Name:       row.Value,
			Antecedent: row.Antecedent,
			PostCount:  row.PostCount,
			Link:       apiHostHTTPS + "/posts?tags=" + url.QueryEscape(row.Value),
		})
	}
	return result, nil
}

// Lookup a list of tags
func (c *client) GetTags(ctx context.Context, tags []string) (map[string]*Tag, error) {
	ctx, cancel := c.httpContext(ctx)
	defer cancel()

	var result []*Tag
	builder := c.tagsResource().
		Param("search[name_comma]", commaJoin(tags)).
		ToJSON(&result)

	if err := c.fetch(ctx, builder); err != nil {
		return nil, err
	}

	return indexTagsByName(ctx, result), nil
}

func indexTagsByName(ctx context.Context, tags []*Tag) map[string]*Tag {
	// Reduce array of results into mapping by tag name
	mapping := make(map[string]*Tag)
	for _, tag := range tags {
		if seen, ok := mapping[tag.Name]; ok {
			// On collision, ignore if same ID was encountered
			if tag.ID == seen.ID {
				continue
			}
			// Otherwise, keep the already-seen tag and discard the doppelganger
			// AFAIK this shouldn't happen, but who knows
			log.Infof(
				ctx,
				"collision of tag name '%s', already seen %d, now discarding %d",
				seen.Name, seen.ID, tag.ID,
			)
		}
		mapping[tag.Name] = tag
	}
	return mapping
}

// GetTagsMatching returns tags that match a particular search pattern.
// Use asterisk (*) as a wildcard in the search pattern.
// https://danbooru.donmai.us/wiki_pages/api%3Aposts
func (c *client) GetTagsMatching(ctx context.Context, pattern string) ([]*Tag, error) {
	ctx, cancel := c.httpContext(ctx)
	defer cancel()

	var result []*Tag
	builder := c.tagsResource().
		Param("search[name_matches]", pattern).
		ToJSON(&result)

	err := c.fetch(ctx, builder)
	return result, err
}
