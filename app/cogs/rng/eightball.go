package rng

import (
	"context"
	"fmt"

	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/fiffu/arisa3/lib/functional"
)

const (
	EightBallQuestion = "a_burning_question"
)

var eightBallResponses = []string{
	"It is certain.",
	"It is decidedly so.",
	"Without a doubt",
	"Yes - definitely",
	"You may rely on it.",
	"As I see it, yes.",
	"Most likely.",
	"Outlook good.",
	"Yes.",
	"Signs point to yes.",
	"It is inevitable.",
	"The market demands it.",

	"Vision unclear, try again.",
	"Ask again later.",
	"Better not tell you now.",
	"Cannot predict now.",
	"Concentrate and ask again.",
	"I'm not legally allowed to comment on that.",

	"Don't count on it.",
	"No.",
	"My sources say no.",
	"Outlook not so good.",
	"Very doubtful.",
	"My calculations say no.",
}

func (c *Cog) eightBallCommand() *types.Command {
	return types.NewCommand("8ball").ForChat().
		Desc("Concentrate and ask again").
		Options(
			types.NewOption(EightBallQuestion).
				Desc("(yes/no questions work well)").
				String().Required(),
		).
		Handler(c.eightBall)
}

func (c *Cog) eightBall(ctx context.Context, req types.ICommandEvent) error {
	asker := formatAsker(req)
	question, _ := req.Args().String(EightBallQuestion)
	reply := functional.TakeRandom(eightBallResponses)

	embed := types.NewEmbed().Description(reply)

	title := fmt.Sprintf("%s: %s", asker, utils.PrettifyCustomEmoji(question))
	msg := fmt.Sprintf("**%s**", reply)
	embed.Description(title + "\n\n" + msg)

	resp := types.NewResponse().Embeds(embed)
	return req.Respond(ctx, resp)
}
