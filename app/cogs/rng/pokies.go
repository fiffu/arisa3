package rng

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/rs/zerolog/log"
)

const (
	PokiesSize        = "grid_size"
	pokiesDefaultRows = 1
	pokiesDefaultCols = 3
)

type cachedEmojis struct {
	guildID string
	emojis  []*discordgo.Emoji
}

func (ce *cachedEmojis) CacheKey() string             { return ce.guildID }
func (ce *cachedEmojis) CacheData() interface{}       { return ce.emojis }
func (ce *cachedEmojis) CacheDuration() time.Duration { return 1 * time.Hour }

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

func (c *Cog) pokies(req types.ICommandEvent) error {
	reply, err := c.getReply(req)
	if err != nil {
		return err
	}
	resp := types.NewResponse().Content(reply)
	return req.Respond(resp)
}

func (c *Cog) getReply(req types.ICommandEvent) (string, error) {
	// Query emoji palette
	guildID := req.Interaction().GuildID
	if guildID == "" {
		return "You need to be in a guild for this command to work!", nil
	}

	emojis, err := c.pullEmojis(req, guildID)
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

func (c *Cog) pullEmojis(req types.ICommandEvent, guildID string) ([]*discordgo.Emoji, error) {
	// Cache lookup
	if cached, ok := c.pokiesCache.Peek(guildID); ok {
		if mem, ok := (cached.CacheData()).([]*discordgo.Emoji); ok {
			return mem, nil
		} else {
			return nil, errors.New("error coercing cached emojis")
		}
	}
	engine.CommandLog(c, req, log.Info()).Msgf("Cache miss")

	emojis, err := req.Session().GuildEmojis(guildID)
	c.pokiesCache.Put(&cachedEmojis{
		guildID,
		emojis,
	})

	engine.CogLog(c, log.Info()).Msgf(
		"Pulled %d emojis from guild id='%s', err=%v",
		len(emojis), guildID, err,
	)

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

	result = strings.Join(grid, "\n")
	sizeCheck = true
	if len(result) > utils.MAX_MESSAGE_LENGTH {
		return "", false
	}
	return
}
