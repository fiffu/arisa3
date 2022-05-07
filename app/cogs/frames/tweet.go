package frames

import "github.com/fiffu/arisa3/app/types"

func (c *Cog) tweetCommand() *types.Command {
	return types.NewCommand("tweet").ForChat().
		Desc("Tweet something, probably").
		Handler(c.tweet)
}

func (c *Cog) tweet(req types.ICommandEvent) error {
	// TODO
	return nil
}
