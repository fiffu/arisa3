package cardboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
)

const embedColour = 0xa4815e

func (d *domain) formatZeroResults(q IQueryPosts) types.IEmbed {
	return types.NewEmbed().
		Description(fmt.Sprintf(
			"I couldn't find any results with `%s`. Try something else?",
			q.Term(),
		))
}

func (d *domain) formatResult(ctx context.Context, query IQueryPosts, posts []*api.Post) (types.IEmbed, error) {
	if len(posts) == 0 {
		return d.formatZeroResults(query), nil
	}
	post := posts[0]

	log.Infof(ctx, "Generating embed for post md5=%s", post.MD5)

	term := query.Term()
	artistTags := splitTags(post.ArtistTags)
	copyrightTags := splitTags(post.CopyrightTags)

	tagsToFetch := []string{term}
	tagsToFetch = append(tagsToFetch, artistTags...)
	tagsToFetch = append(tagsToFetch, copyrightTags...)

	tagData, err := d.client.GetTags(ctx, tagsToFetch)
	if err != nil {
		tagData = make(map[string]*api.Tag)

		log.Infof(ctx, "ignored error fetching tag data for post md5=%s, err=%s", post.MD5, err)
	}

	title := embedTitle(post)
	image := post.GetFileURL()

	artistsField := embedFieldTags(artistTags, tagData)
	sourcesField := embedFieldTags(copyrightTags, tagData)
	linksField := d.embedLinks(query, post)

	footer := fmt.Sprintf("Matched against: " + term)
	if termTag, ok := tagData[term]; ok {
		footer += fmt.Sprintf(" (%d)", termTag.PostCount)
	}

	inline := true
	embed := types.NewEmbed().
		Colour(embedColour).
		Title(title).
		Image(image).
		Footer(footer, "")

	// Conditionally append fields
	if artistsField != "" {
		embed.Field("Artist", artistsField, inline)
	}
	if sourcesField != "" {
		embed.Field("Source", sourcesField, inline)
	}
	// Links field always comes last
	embed.Field("Links", linksField, !inline)

	return embed, nil
}

func embedTitle(post *api.Post) string {
	maxLen := 255

	artists := embedTitleArtists(post)
	if len(artists) > maxLen/2 {
		artists = ""
	}

	if title, ok := fitString(
		parseTags(post.CharacterTags),
		", ",
		" and ",
		" and %d more",
		artists,
		maxLen,
	); ok {
		if title != "" {
			return title
		}
	}

	if len(artists) < maxLen && artists != "" {
		if artists[0] == ' ' {
			artists = artists[1:]
		}
		return artists
	}

	return fmt.Sprint("Picture #", post.ID)
}

func embedTitleArtists(post *api.Post) string {
	tags := parseTags(post.ArtistTags)
	artists := ""
	switch len(tags) {
	case 0:
		artists = ""
	case 1:
		artists = " drawn by " + tags[0]
	default:
		artists = " drawn by " + tags[0] + " and others"
	}
	return artists
}

func embedFieldTags(tagNames []string, tagData map[string]*api.Tag) string {
	if len(tagNames) == 0 {
		return "(none)"
	}

	for i, name := range tagNames {
		url := api.GetSearchURL(name)
		data, ok := tagData[name]
		if !ok {
			tagNames[i] = fmt.Sprintf("[`%s`](%s)", name, url)
		} else {
			count := data.PostCount
			tagNames[i] = fmt.Sprintf("[`%s`](%s) (%d)", name, url, count)
		}
	}

	maxLen := utils.MAX_MESSAGE_LENGTH / 2
	if res, ok := fitString(
		tagNames,
		"\n",
		"\n",
		"\nand %d more",
		"",
		maxLen,
	); ok {
		return res
	} else {
		return "(too many)"
	}
}

func (d *domain) embedLinks(query IQueryPosts, post *api.Post) string {
	queryStr := strings.Join(query.Tags(), " ")
	return fmt.Sprintf(
		"[pic](%s) · [post](%s) · [search](%s)",
		post.GetFileURL(),
		api.GetPostURL(post),
		api.GetSearchURL(queryStr),
	)
}

func splitTags(str string) []string {
	if str == "" {
		return []string{}
	}
	tags := strings.Split(str, " ")
	return tags
}

func parseTags(str string) []string {
	tags := splitTags(str)
	for i, tag := range tags {
		tags[i] = utils.EscapeMarkdown(untaggify(tag))
	}
	return tags
}

func fitString(
	strs []string,
	sep, sepLast, sepOverflowf, mustAppend string,
	totalLength int) (string, bool) {

	if len(strs) == 0 {
		return "", true
	}

	if totalLength < 0 {
		totalLength = 9999999
	}

	joined := joinWithTail(strs, sep, sepLast)
	tailCount := 0

	for len(joined+mustAppend) > totalLength {
		tailCount += 1
		head := strs[:len(strs)-tailCount]
		joined = joinWithTailf(head, sep, sepOverflowf, tailCount)
		if len(head) == 0 {
			return "", false
		}
	}
	return joined + mustAppend, true
}

func joinWithTail(strs []string, joiner, penult string) string {
	if len(strs) == 1 {
		return strs[0]
	}

	lastIdx := len(strs) - 1
	last := strs[lastIdx]
	joined := strings.Join(strs, joiner)
	return joined + penult + last
}

func joinWithTailf(strs []string, joiner, tailFmt string, args ...interface{}) string {
	joined := strings.Join(strs, joiner)
	return joined + fmt.Sprintf(tailFmt, args...)
}
