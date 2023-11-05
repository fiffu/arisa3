package cardboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/lib/functional"
)

func (c *Cog) tagSuggestCommand() *types.Command {
	return types.NewCommand("tags").ForChat().
		Desc("See suggested tags that match the given query.").
		Options(
			types.NewOption(OptionQuery).
				String().Required(),
		).
		Handler(c.tagAutocompleteSuggest)
}

func (c *Cog) tagAutocompleteSuggest(ctx context.Context, req types.ICommandEvent) error {
	queryStr, _ := req.Args().String(OptionQuery)

	suggestedTags, err := c.domain.TagsSearch(queryStr)
	if err != nil {
		return err
	}

	var desc string
	if len(suggestedTags) > 0 {
		lines := functional.Map(suggestedTags, func(tag *api.TagSuggestion) string {
			return fmt.Sprintf("%s (%d posts)", tag.Name, tag.PostCount)
		})
		desc = strings.Join(lines, "\n")
	} else {
		desc = fmt.Sprintf("There's no tags that match '%s'", desc)
	}

	emb := types.NewEmbed().
		Colour(embedColour).
		Title(fmt.Sprintf("Tags matching '%s'", queryStr)).
		Description(desc)
	resp := types.NewResponse().Embeds(emb)
	return req.Respond(ctx, resp)
}
