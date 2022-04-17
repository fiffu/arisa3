package rng

import (
	"arisa3/app/types"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// Command consts
const (
	RollCommand       = "roll"
	RollOptionExpr    = "expression"
	RollOptionComment = "comment"
)

// Regex patterns
var (
	_20      *regexp.Regexp = regexp.MustCompile(`^\d+$`)                 // "20"
	_d20     *regexp.Regexp = regexp.MustCompile(`^[dD]\d+$`)             // "d20"
	_d20mod  *regexp.Regexp = regexp.MustCompile(`^[dD]\d+[\+\-]\d+$`)    // "d20+5"
	_3d20    *regexp.Regexp = regexp.MustCompile(`^\d+[dD]\d+$`)          // "3d20"
	_3d20mod *regexp.Regexp = regexp.MustCompile(`^\d+[dD]\d+[\+\-]\d+$`) // "3d20+5"
)

// dice describes the dice to be thrown
type dice struct {
	// Number of dice in the roll
	count int
	// Sides per die in the roll
	sides int
	// Arithmetic modifier (summed after the result of rolling dice*sides)
	mod int

	// Whether this value was parsed successfully from user
	parsed bool
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
	var expression, comment string
	if value, ok := req.Args().String("expression"); ok {
		expression = value
	}
	if value, ok := req.Args().String("comment"); ok {
		comment = value
	}

	d := c.parseExpr(req, expression)

	// if expression couldn't be parsed, treat it as a comment
	if !d.parsed && comment == "" {
		comment = expression
	}
	result := toss(d)
	resp := formatResponse(req, d, result, comment)
	return req.Respond(resp)
}

func (c *Cog) parseExpr(req types.ICommandEvent, s string) dice {
	s = strings.ToLower(s)

	var count, sides, mod string
	switch {

	case _20.MatchString(s):
		sides = s
		return dice{
			count:  1,
			sides:  Atoi(sides),
			mod:    0,
			parsed: true,
		}

	case _d20.MatchString(s):
		_, sides = SplitOnce(s, "d")
		return dice{
			count:  1,
			sides:  Atoi(sides),
			mod:    0,
			parsed: true,
		}

	case _d20mod.MatchString(s):
		_, s = SplitOnce(s, "d")
		var sides, mod string
		if strings.Contains(s, "+") {
			sides, mod = SplitOnce(s, "+")
		} else {
			sides, mod = SplitOnce(s, "-")
			mod = "-" + mod
		}
		return dice{
			count:  1,
			sides:  Atoi(sides),
			mod:    Atoi(mod),
			parsed: true,
		}

	case _3d20.MatchString(s):
		count, sides = SplitOnce(s, "d")
		return dice{
			count:  Atoi(count),
			sides:  Atoi(sides),
			mod:    0,
			parsed: true,
		}

	case _3d20mod.MatchString(s):
		count, s = SplitOnce(s, "d")
		if strings.Contains(s, "+") {
			sides, mod = SplitOnce(s, "+")
		} else {
			sides, mod = SplitOnce(s, "-")
			mod = "-" + mod
		}
		return dice{
			count:  Atoi(count),
			sides:  Atoi(sides),
			mod:    Atoi(mod),
			parsed: true,
		}

	default:
		return dice{
			count:  1,
			sides:  100,
			mod:    -1,
			parsed: false,
		}
	}
}

func Atoi(s string) int {
	if num, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return num
	}
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
	for i := 0; i < count; i++ {
		sum += throwDie(d.sides)
	}
	sum += d.mod
	return
}

func throwDie(sides int) int {
	if sides <= 0 {
		return 0
	}
	return 1 + rand.Intn(sides)
}

func formatResponse(req types.ICommandEvent, d dice, result int, comment string) types.ICommandResponse {
	whatDice := "0-99"
	resultStr := fmt.Sprintf("%2d", result)
	if d.parsed {
		whatDice = formatDice(d)
		resultStr = fmt.Sprintf("%d", result)
	}
	msg := fmt.Sprintf("Rolling %s: **%s**\nComment: %s\nDice: %+v", whatDice, resultStr, comment, d)
	return types.NewResponse().Content(msg)
}

func formatDice(d dice) string {
	if !d.parsed {
		return "0-99"
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
