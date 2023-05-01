package colours

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"strings"
	"time"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/fiffu/arisa3/lib/functional"
	"github.com/rs/zerolog/log"
)

func (c *Cog) colInfoCommand() *types.Command {
	return types.NewCommand("colinfo").ForChat().
		Desc("Tells you about your colour").
		Handler(func(req types.ICommandEvent) error {
			return c.colInfo(req)
		})
}

func (c *Cog) colInfo(req types.ICommandEvent) error {
	mem, resp, err := c.fetchMember(req)
	if err != nil {
		return err
	}
	if resp != nil {
		return req.Respond(resp)
	}

	guildID := mem.Guild().ID()
	userID := mem.UserID()

	role := c.domain.GetColourRole(mem)
	if role == nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("No colour role found, guild=%s user=%s", guildID, userID)
		return req.Respond(types.NewResponse().
			Content("You don't have a colour role. Use /col to get a random colour!"))
	}

	rerollCDEndTime, err := c.domain.GetRerollCooldownEndTime(mem)
	if err != nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting cooldown end time, guild=%s user=%s", guildID, userID)
		return err
	}

	lastMutateTime, _, err := c.domain.GetLastMutate(mem)
	if err != nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting last mutate time, guild=%s user=%s", guildID, userID)
		return err
	}

	lastFrozenTime, err := c.domain.GetLastFrozen(mem)
	if err != nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting last frozen time, guild=%s user=%s", guildID, userID)
		return err
	}

	history, err := c.domain.GetHistory(mem)
	if err != nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored getting last frozen time, guild=%s user=%s", guildID, userID)
		return err
	}
	engine.CommandLog(c, req, log.Info()).
		Msgf("Colour history guild=%s user=%s: %v", guildID, userID, functional.Map(
			history.records, func(c *ColoursLogRecord) string { return c.ColourHex },
		))

	info, err := c.formatColInfo(time.Now(), rerollCDEndTime, lastMutateTime, lastFrozenTime, history)
	if err != nil {
		engine.CommandLog(c, req, log.Error()).Err(err).
			Msgf("Errored formatting colour info, guild=%s user=%s", guildID, userID)
		return err
	}

	reply := types.NewResponse()
	embed := newEmbed(role.Colour()).Description(info.desc)
	if info.img.ok {
		img := info.img
		reply.File(img.filename, img.contentType, img.file)
		embed.Image("attachment://" + img.filename)
	}

	return req.Respond(reply.Embeds(embed))
}

type colInfo struct {
	desc string
	img  struct {
		ok                               bool
		file                             *bytes.Buffer
		filename, extension, contentType string
	}
}

func (c *Cog) formatColInfo(
	now time.Time,
	rerollCDEndTime, lastMutateTime, lastFrozenTime time.Time,
	history *History,
) (*colInfo, error) {
	desc := make([]string, 0)

	desc = append(desc, "**Reroll cooldown:**")
	if now.Before(rerollCDEndTime) {
		desc = append(desc, utils.FormatDuration(rerollCDEndTime.Sub(now)))
	} else {
		desc = append(desc, "_(No cooldown, reroll available)_")
	}

	desc = append(desc, "", "**Last mutate:**")
	if lastMutateTime == Never {
		desc = append(desc, "_(Never)_")
	} else if now.After(lastMutateTime) {
		diff := now.Sub(lastMutateTime)
		desc = append(desc, utils.FormatDuration(diff)+" ago")
	} else {
		desc = append(desc, "Moments ago")
	}

	if lastFrozenTime != Never {
		frozenDuration := utils.FormatDuration(now.Sub(lastFrozenTime))
		desc = append(desc, "Frozen "+frozenDuration+" ago")
	}

	ret := &colInfo{}

	if history != nil && len(history.records) > 0 {
		buf, ext, mime, err := formatColHistory(history, time.Duration(c.cfg.MutateCooldownMins)*time.Minute)
		if err != nil {
			return nil, err
		}

		desc = append(desc, "", "**Image history, newest → oldest:**")
		ret.img.ok = true
		ret.img.file = buf
		ret.img.filename = "history." + ext
		ret.img.extension = ext
		ret.img.contentType = mime
	}

	ret.desc = strings.Join(desc, "\n")
	return ret, nil
}

func formatColHistory(h *History, interval time.Duration) (file *bytes.Buffer, fileExt, fileContent string, err error) {
	colours := partitionColours(h, interval)
	pixelsPerInterval := 4

	buf := bytes.NewBuffer(make([]byte, 0))
	err = png.Encode(buf, horizontalPartitionImage{
		partitions:      colours,
		partitionWidth:  pixelsPerInterval,
		partitionHeight: pixelsPerInterval * 5,
	})
	file = bytes.NewBuffer(buf.Bytes())
	fileExt = "png"
	fileContent = "image/png"
	return file, fileExt, fileContent, err
}

type horizontalPartitionImage struct {
	partitions      []*Colour
	partitionWidth  int
	partitionHeight int
}

func (b horizontalPartitionImage) ColorModel() color.Model { return color.RGBAModel }
func (b horizontalPartitionImage) At(x, y int) color.Color { return b.partitions[x/b.partitionWidth] }
func (b horizontalPartitionImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{
			X: b.partitionWidth * len(b.partitions),
			Y: b.partitionHeight,
		},
	}
}

func partitionColours(h *History, interval time.Duration) []*Colour {
	if len(h.records) == 0 {
		return make([]*Colour, 0)
	}

	historySpan := h.end.Sub(h.start)
	numSpans := int(math.Ceil(historySpan.Seconds() / interval.Seconds()))

	rec := h.records[0]
	recs := h.records[1:]
	spans := make([]*Colour, numSpans)

	for retIdx := range spans {
		spanEnd := h.start.Add(time.Duration(retIdx+1) * interval)

		for len(recs) > 0 && recs[0].TStamp.Unix() < spanEnd.Unix() {
			rec = recs[0]
			recs = recs[1:]
		}

		spans[retIdx] = (&Colour{}).FromRGBHex(rec.ColourHex)
	}
	return spans
}
