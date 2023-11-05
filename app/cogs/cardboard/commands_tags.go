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

	suggestedTags, err := c.domain.TagsSearch(ctx, queryStr)
	if err != nil {
		return err
	}

	var desc string
	if len(suggestedTags) > 0 {
		lines := functional.Map(suggestedTags, formatSuggestion)
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

func formatSuggestion(suggest *api.TagSuggestion) string {
	var ante string
	if suggest.Antecedent != "" {
		ante = fmt.Sprintf(" ← _alias from '%s'_", suggest.Antecedent)
	}
	return fmt.Sprintf("[`%s`](%s) (%d)%s", suggest.Name, suggest.Link, suggest.PostCount, ante)
}
