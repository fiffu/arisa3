package api

import "github.com/rs/zerolog/log"

type Tag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	PostCount int    `json:"post_count"`
}

// Lookup a list of tags
func (c *client) GetTags(tags []string) (map[string]*Tag, error) {
	ctx, cancel := c.httpContext()
	defer cancel()

	var result []*Tag
	builder := c.requestTags().
		Param("search[name_comma]", commaJoin(tags)).
		ToJSON(&result)

	if err := c.fetch(ctx, builder); err != nil {
		return nil, err
	}

	return indexTagsByName(result), nil
}

func indexTagsByName(tags []*Tag) map[string]*Tag {
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
			log.Warn().Msgf(
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
func (c *client) GetTagsMatching(pattern string) ([]*Tag, error) {
	ctx, cancel := c.httpContext()
	defer cancel()

	var result []*Tag
	builder := c.requestTags().
		Param("search[name_matches]", pattern).
		ToJSON(&result)

	err := c.fetch(ctx, builder)
	return result, err
}
