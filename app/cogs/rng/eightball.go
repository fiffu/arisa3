package rng

import (
	"arisa3/app/types"
	"math/rand"
)

const (
	EightBallCommand = "8ball"
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
	return types.NewCommand(EightBallCommand).ForChat().
		Desc("Concentrate and ask again").
		Handler(c.eightBall)
}

func (c *Cog) eightBall(req types.ICommandEvent) error {
	reply := randChoice(eightBallResponses)
	return req.Respond(
		types.NewResponse().Content(reply),
	)
}

func randChoice(slc []string) string {
	size := len(slc)
	n := rand.Intn(size)
	return slc[n]
}
