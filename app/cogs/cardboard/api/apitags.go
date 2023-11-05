package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

// AutocompleteTag implements api.IClient.
// Note: This is not part of the Danbooru API proper.
// It is undocumented, and returns a HTML response.
func (c *client) AutocompleteTag(ctx context.Context, query string) ([]*TagSuggestion, error) {
	ctx, cancel := c.httpContext(ctx)
	defer cancel()

	buf := new(strings.Builder)
	builder := c.autocompleteEndpoint().
		Param("search[query]", query).
		Param("search[type]", "tag_query").
		Param("version", "1").
		Param("limit", "10").
		ToWriter(buf)
	if err := c.fetch(ctx, builder); err != nil {
		return nil, err
	}

	reader := strings.NewReader(buf.String())
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}

	var result []*TagSuggestion
	const selector = "li.ui-menu-item"
	doc.Find(selector).
		Each(func(idx int, s *goquery.Selection) {
			suggestion, err := parseAutocompleteElem(ctx, s)
			if err != nil {
				log.Infof(ctx, "Failed to parse autocomplete elem at %s[%d], err: %v", selector, idx, err)
			}

			result = append(result, suggestion)
		})
	return result, nil
}

func parseAutocompleteElem(ctx context.Context, s *goquery.Selection) (*TagSuggestion, error) {
	suggest := &TagSuggestion{}

	tagName, ok := s.Attr("data-autocomplete-value")
	if !ok {
		return nil, fmt.Errorf("missing attr: 'data-autocomplete-value'")
	}
	suggest.Name = tagName

	if span := s.Find("span.post-count").First(); span != nil {
		rawPostCount := span.Text()
		if postCount, err := strconv.Atoi(rawPostCount); err != nil {
			return nil, fmt.Errorf("failed to parse attr 'span.post-count', err: %w", err)
		} else {
			suggest.PostCount = postCount
		}
	}

	if ante := s.Find("span.autocomplete-antecedent").First(); ante != nil {
		suggest.Antecedent = ante.Text()
	}

	if a := s.Find("a"); a != nil {
		if href, ok := a.Attr("href"); ok {
			suggest.Link = apiHostHTTPS + href
		}
	}

	return suggest, nil
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
