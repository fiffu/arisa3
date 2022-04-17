package rng

import (
	"arisa3/app/engine"
	"arisa3/app/types"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// Command consts
const (
	RollCommand       = "roll"
	RollOptionExpr    = "expression"
	RollOptionComment = "comment"
)

// Regex patterns
var validRollPatterns = []*regexp.Regexp{
	regexp.MustCompile(`^\d+$`),                 // "20"
	regexp.MustCompile(`^[dD]\d+$`),             // "d20"
	regexp.MustCompile(`^[dD]\d+[\+\-]\d+$`),    // "d20+5"
	regexp.MustCompile(`^\d+[dD]\d+$`),          // "3d20"
	regexp.MustCompile(`^\d+[dD]\d+[\+\-]\d+$`), // "3d20+5"
}

// dice describes the dice to be thrown
type dice struct {
	// Number of dice in the roll
	count int
	// Sides per die in the roll
	sides int
	// Arithmetic modifier (summed after the result of rolling dice*sides)
	mod int
	// User's comment
	comment string

	// Whether this value was parsed successfully from user
	parsed bool
}

func (r dice) NewWithDefaults() dice {
	return dice{
		// default roll [0-99] i.e. d100-1
		count: 1,
		sides: 100,
		mod:   -1,
	}
}

func (c *Cog) rollCommand() *types.Command {
	return types.NewCommand(RollCommand).ForChat().
		Desc("Rolls dice (supports algebraic notation)").
		Options(
			types.NewOption(RollOptionExpr).Desc("dice notation, such as *3d5+10*").
				String(),
			types.NewOption(RollOptionComment).Desc("optional comment").
				String(),
		).
		Handler(c.roll)
}

func (c *Cog) roll(req types.ICommandEvent) error {
	// resp := types.NewResponse()

	d := dice{}.NewWithDefaults()

	var expression, comment string
	if value, ok := req.Args().String("expression"); ok {
		expression = value
	}
	if value, ok := req.Args().String("comment"); ok {
		comment = value
	}

	d = c.parseExpr(req, expression)

	// if parsing failed, we treat the `expression` argument as a comment
	if !d.parsed && comment == "" {
		comment = expression
	}
	result := toss(d)
	resp := formatResponse(req, d, comment, result)
	return req.Respond(resp)
}

func (c *Cog) parseExpr(req types.ICommandEvent, s string) dice {
	val := dice{}.NewWithDefaults()
	count, sides, mod, ok := findMatches(s)
	if !ok {
		return val
	}
	val.parsed = true

	for fld, v := range map[string]string{"c": count, "s": sides, "m": mod} {
		num, err := strconv.Atoi(v)
		if err != nil {
			engine.CommandLog(c, req, log.Error()).Err(err).
				Msgf("failed converting '%v' to int for %s", v, fld)
			continue
		}
		// assign
		switch fld {
		case "c":
			val.count = num
		case "s":
			val.sides = num
		case "m":
			val.mod = num
		}
	}

	if comment, ok := req.Args().String(RollOptionComment); ok {
		val.comment = comment
	}
	return val
}

func findMatches(s string) (count, sides, mod string, matched bool) {
	s = strings.ToLower(s)

	for _, patt := range validRollPatterns {
		matched = patt.MatchString(s)
	}
	if !matched {
		return
	}

	// example: "3d12+15"
	if strings.Contains(s, "d") {
		count, s = SplitOnce(s, "d")
	}
	// remaining is "12+15"
	switch {
	case strings.Contains(s, "+"):
		sides, mod = SplitOnce(s, "+")
	case strings.Contains(s, "-"):
		sides, mod = SplitOnce(s, "-")
	default:
		sides = s
	}

	return
}

func SplitOnce(s, substr string) (left, right string) {
	pivot := strings.Index(s, substr)
	offset := pivot + len(substr)
	left = s[:pivot]
	right = s[offset:]
	return
}

func toss(d dice) (sum int) {
	count := d.count
	lo := 1
	hi := d.sides
	for i := 0; i < count; i++ {
		sum += throwDie(lo, hi)
	}
	return
}

func throwDie(lo, hi int) int {
	return rand.Intn(lo+hi) - lo
}

func formatResponse(req types.ICommandEvent, d dice, comment string, sum int) types.ICommandResponse {
	whatDice := formatDice(d)
	msg := fmt.Sprintf("Rolling %s: **%d**\nComment: %s", whatDice, sum, comment)
	return types.NewResponse().Content(msg)
}

func formatDice(d dice) string {
	if !d.parsed {
		return "1-99"
	}
	var count, sides, mod, modSign string
	if d.count != 1 {
		count = fmt.Sprintf("%d", d.count)
	}
	sides = fmt.Sprintf("%d", d.sides)
	switch {
	case d.mod == 0:
		mod = ""
		modSign = ""
	case d.mod > 0:
		mod = fmt.Sprintf("%d", d.mod)
		modSign = "+"
	case d.mod < 0:
		mod = fmt.Sprintf("-%d", d.mod)
		modSign = "-"
	}
	return count + "d" + sides + modSign + mod
}
