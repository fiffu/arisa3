package cardboard

import "github.com/fiffu/arisa3/app/types"

func (c *Cog) tagSuggestCommand() *types.Command {
	return types.NewCommand("tags").ForChat().
		Desc("See suggested tags that match the given query.").
		Options(
			types.NewOption(OptionQuery).
				String().Required(),
		).
		Handler(c.promote)
}
