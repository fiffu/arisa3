package rng

import (
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/rs/zerolog/log"
)

const (
	PokiesSize = "grid_size"
)

func (c *Cog) pokiesCommand() *types.Command {
	return types.NewCommand("pokies").ForChat().
		Desc("An avian slots machine. Use with care.").
		Options(
			types.NewOption(PokiesSize).
				Int().
				Desc("generate result in N x N grid; if not given, result will be a single row of 3 x 1"),
		).
		Handler(c.pokies)
}

func (c *Cog) pokies(req types.ICommandEvent) error {
	// Query emoji palette
	guildID := req.Interaction().GuildID
	emojis, err := req.Session().GuildEmojis(guildID)
	if err != nil {
		return err
	}

	// Parse request, build reply
	rows, cols, tooBig := parseGrid(req)
	var reply string
	if tooBig {
		reply = "That's just way too much work " + utils.BIRB
	} else {
		reply = buildGrid(rows, cols, emojis)
	}

	resp := types.NewResponse().Content(reply)
	return req.Respond(resp)
}

func parseGrid(req types.ICommandEvent) (rows int, cols int, tooBig bool) {
	rows, cols = 1, 3
	if size, ok := req.Args().Int(PokiesSize); ok {
		if size > 9 {
			tooBig = true
			return
		}
		rows = size
		cols = size
	}
	return
}

func buildGrid(rows, cols int, emojis []*discordgo.Emoji) string {
	log.Info().Msgf("Building grid of %d x %d", cols, rows)

	// Reset the rand seed otherwise it will always yield the same result
	rand.Seed(time.Now().Unix())

	emojiCount := len(emojis)
	grid := make([]string, rows)

	for y := range grid {
		row := make([]string, cols)
		for x := range row {
			n := rand.Intn(emojiCount)
			row[x] = emojis[n].MessageFormat()
		}
		grid[y] = strings.Join(row, " ")
	}

	return strings.Join(grid, "\n")
}
