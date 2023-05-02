package rng

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/fiffu/arisa3/lib/functional"
)

const (
	PokiesSize        = "grid_size"
	pokiesDefaultRows = 1
	pokiesDefaultCols = 3
	pokiesNumSymbols  = 20 // Number of symbols per slot
)

type cachedEmojis struct {
	guildID string
	emojis  []*discordgo.Emoji
}

func (ce *cachedEmojis) CacheKey() string { return ce.guildID }

func (c *Cog) pokiesCommand() *types.Command {
	return types.NewCommand("pokies").ForChat().
		Desc("An avian slots machine. Use with care.").
		Options(
			types.NewOption(PokiesSize).
				Desc("generate result in N x N grid; if not given, result will be a single row of 3 x 1").
				Int(),
		).
		Handler(c.pokies)
}

func (c *Cog) pokies(ctx context.Context, req types.ICommandEvent) error {
	reply, err := c.getReply(ctx, req)
	if err != nil {
		return err
	}
	resp := types.NewResponse().Content(reply)
	return req.Respond(ctx, resp)
}

func (c *Cog) getReply(ctx context.Context, req types.ICommandEvent) (string, error) {
	// Query emoji palette
	guildID := req.Interaction().GuildID
	if guildID == "" {
		return "You need to be in a guild for this command to work!", nil
	}

	emojis, err := c.pullEmojis(ctx, req, guildID)
	if err != nil {
		return "", err
	}

	// Parse request, build reply
	replyTooBig := "That's just way too much work " + utils.BIRB
	rows, cols, ok := parseGrid(req)
	if !ok {
		return replyTooBig, nil
	}
	result, ok := buildGrid(rows, cols, emojis)
	if !ok {
		return replyTooBig, nil
	}
	return result, nil
}

func (c *Cog) pullEmojis(ctx context.Context, req types.ICommandEvent, guildID string) ([]*discordgo.Emoji, error) {
	// Cache lookup
	if cached, ok := c.pokiesCache.Peek(guildID); ok {
		return cached.emojis, nil
	}
	log.Debugf(ctx, "Cache miss")

	emojis, err := req.Session().GuildEmojis(guildID)
	c.pokiesCache.Put(&cachedEmojis{
		guildID,
		emojis,
	})

	log.Infof(ctx, "Pulled %d emojis from guild id='%s', err=%v", len(emojis), guildID, err)
	return emojis, err
}

func parseGrid(req types.ICommandEvent) (rows int, cols int, sizeCheck bool) {
	sizeCheck = true
	rows, cols = pokiesDefaultRows, pokiesDefaultCols
	if size, ok := req.Args().Int(PokiesSize); ok {
		if size > 8 {
			sizeCheck = false
			return
		}
		if size > 0 {
			rows = size
			cols = size
		}
	}
	return
}

func buildGrid(rows, cols int, emojis []*discordgo.Emoji) (result string, sizeCheck bool) {
	// Reset the rand seed otherwise it will always yield the same result
	rand.Seed(time.Now().Unix())

	slotPool := functional.SliceOf(emojis).
		Shuffle().
		Take(pokiesNumSymbols)

	grid := make([]string, rows)

	for y := range grid {
		row := make([]string, cols)
		for x := range row {
			row[x] = slotPool.TakeRandom().MessageFormat()
		}
		grid[y] = strings.Join(row, " ")
	}

	result = strings.Join(grid, "\n")
	sizeCheck = true
	if len(result) > utils.MAX_MESSAGE_LENGTH {
		return "", false
	}
	return
}
