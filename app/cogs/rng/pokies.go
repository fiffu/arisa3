package rng

import (
	"math/rand"
	"strings"
	"time"

	"github.com/fiffu/arisa3/app/types"
	"github.com/rs/zerolog/log"
)

const (
	PokiesSize = "n"
)

func (c *Cog) pokiesCommand() *types.Command {
	return types.NewCommand("pokies").ForChat().
		Desc("An avian slots machine. Use with care.").
		Options(
			types.NewOption(PokiesSize).
				Int().
				Desc("generate result in NxN grid; if not given, result will be a single 3x1 row"),
		).
		Handler(c.pokies)
}

func (c *Cog) pokies(req types.ICommandEvent) error {
	// Reset the rand seed
	rand.Seed(time.Now().Unix())

	guildID := "118320695057842183" // req.Interaction().GuildID

	emojis, err := req.Session().GuildEmojis(guildID)
	emojiCount := len(emojis)
	if err != nil {
		return err
	}

	rows, cols := 1, 3
	log.Info().Msgf("args %v", req.Args())
	if size, ok := req.Args().Int(PokiesSize); ok {
		rows = size
		cols = size
	}

	log.Info().Msgf("Building grid of %d x %d", cols, rows)

	grid := make([]string, rows)
	for y := range grid {
		row := make([]string, cols)
		for x := range row {
			n := rand.Intn(emojiCount)
			row[x] = emojis[n].MessageFormat()
		}
		grid[y] = strings.Join(row, " ")
	}

	reply := strings.Join(grid, "\n")
	resp := types.NewResponse().Content(reply)
	return req.Respond(resp)
}
