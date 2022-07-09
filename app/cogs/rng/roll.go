package rng

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/fiffu/arisa3/app/engine"
	"github.com/fiffu/arisa3/app/types"
	"github.com/fiffu/arisa3/app/utils"
	"github.com/fiffu/arisa3/lib"
)

// Command consts
const (
	RollCommand        = "roll"
	RollExpression     = "expression_or_comment"
	DefaultRollContent = "0-99"
)

// Regex patterns
var (
	_20        *regexp.Regexp = regexp.MustCompile(`^\d+$`)              // "20"
	_d20       *regexp.Regexp = regexp.MustCompile(`^D\d+$`)             // "d20"
	_d20modif  *regexp.Regexp = regexp.MustCompile(`^D\d+[\+\-]\d+$`)    // "d20+5"
	_3d20      *regexp.Regexp = regexp.MustCompile(`^\d+D\d+$`)          // "3d20"
	_3d20modif *regexp.Regexp = regexp.MustCompile(`^\d+D\d+[\+\-]\d+$`) // "3d20+5"
)

// dice describes the dice to be thrown
type dice struct {
	// Number of dice in the roll
	count int
	// Sides per die in the roll
	sides int
	// Arithmetic modifier (summed after the result of rolling dice*sides)
	modif int

	// Whether this value was parsed successfully from user
	parsed bool
}

func (c *Cog) rollCommand() *types.Command {
	return types.NewCommand(RollCommand).ForChat().
		Desc("Rolls dice (supports algebraic notation)").
		Options(
			types.NewOption(RollExpression).
				Desc("dice expression like 3d5+10 and/or a comment (default is d100-1)").
				String(),
		).
		Handler(c.roll)
}

func (c *Cog) roll(req types.ICommandEvent) error {
	var input string
	if value, ok := req.Args().String(RollExpression); ok {
		input = value
	}
	if len(input) > 10 {
		return req.Respond(
			types.NewResponse().Content("That's just way too much work " + utils.BIRB),
		)
	}
	d, comment := parse(input)

	if d.sides > 999_999 || d.count > 999_999 {
		return req.Respond(
			types.NewResponse().Content("That's just way too much work " + utils.BIRB),
		)
	}

	// if expression couldn't be parsed, treat it as a comment
	if !d.parsed && comment == "" {
		comment = input
	}
	result := toss(d)
	resp := formatResponse(req, d, result, comment)
	return req.Respond(resp)
}

func parse(input string) (dice, string) {
	delim := " "
	head, comment := lib.SplitOnce(input, delim)
	d := parseExpr(head)
	if !d.parsed && head != "" {
		comment = head + " " + comment
	}
	return d, strings.Trim(comment, " \n\t")
}

func parseExpr(s string) dice {
	s = strings.ToUpper(s)
	s = strings.Trim(s, " \t")

	defaultRoll := dice{
		count:  1,
		sides:  100,
		modif:  -1,
		parsed: false,
	}
	// a bit of optimization...
	idxD := strings.Index(s, "D")
	if len(s) > idxD+1 && !strings.ContainsAny(string(s[idxD+1]), "0123456789") {
		// if has a "d" but next char isn't a number
		return defaultRoll
	}

	var count, sides, modif string
	switch {

	default:
		return defaultRoll

	case _20.MatchString(s):
		return dice{
			count:  1,
			sides:  lib.Atoi(s),
			modif:  0,
			parsed: true,
		}

	case _d20.MatchString(s):
		_, sides = lib.SplitOnce(s, "D")
		return dice{
			count:  1,
			sides:  lib.Atoi(sides),
			modif:  0,
			parsed: true,
		}

	case _d20modif.MatchString(s):
		_, s = lib.SplitOnce(s, "D")
		var sides, modif string
		if strings.Contains(s, "+") {
			sides, modif = lib.SplitOnce(s, "+")
		} else {
			sides, modif = lib.SplitOnce(s, "-")
			modif = "-" + modif
		}
		return dice{
			count:  1,
			sides:  lib.Atoi(sides),
			modif:  lib.Atoi(modif),
			parsed: true,
		}

	case _3d20.MatchString(s):
		count, sides = lib.SplitOnce(s, "D")
		return dice{
			count:  lib.Atoi(count),
			sides:  lib.Atoi(sides),
			modif:  0,
			parsed: true,
		}

	case _3d20modif.MatchString(s):
		count, s = lib.SplitOnce(s, "D")
		if strings.Contains(s, "+") {
			sides, modif = lib.SplitOnce(s, "+")
		} else {
			sides, modif = lib.SplitOnce(s, "-")
			modif = "-" + modif
		}
		return dice{
			count:  lib.Atoi(count),
			sides:  lib.Atoi(sides),
			modif:  lib.Atoi(modif),
			parsed: true,
		}
	}
}

func toss(d dice) (sum int) {
	count := d.count
	for i := 0; i < count; i++ {
		sum += throwDie(d.sides)
	}
	sum += d.modif
	return
}

func throwDie(sides int) int {
	if sides <= 0 {
		return 0
	}
	return 1 + rand.Intn(sides)
}

func formatResponse(req types.ICommandEvent, d dice, result int, comment string) types.ICommandResponse {
	asker := formatAsker(req)
	whatDice := DefaultRollContent
	resultStr := fmt.Sprintf("%2d", result)
	if d.parsed {
		whatDice = formatDice(d)
		resultStr = fmt.Sprintf("%d", result)
	}
	embed := types.NewEmbed().
		Description(fmt.Sprintf("Rolling %s: **%s**", whatDice, resultStr))
	if comment != "" {
		foot := fmt.Sprintf("%s: %s", asker, engine.PrettifyCustomEmoji(comment))
		embed.Footer(foot, "")
	}
	return types.NewResponse().Embeds(embed)
}

func formatAsker(req types.ICommandEvent) string {
	username := req.User().String()

	nickname := ""
	member := req.Interaction().Member
	if member != nil && member.Nick != "" {
		nickname = member.Nick
	}

	if nickname != "" {
		return nickname
	}
	return username
}

func formatDice(d dice) string {
	if !d.parsed {
		return ""
	}
	var count, sides, modif, modSign string
	if d.count != 1 {
		count = fmt.Sprintf("%d", d.count)
	}
	sides = fmt.Sprintf("%d", d.sides)
	switch {
	case d.modif == 0:
		modif = ""
		modSign = ""
	case d.modif > 0:
		modif = fmt.Sprintf("%d", d.modif)
		modSign = "+"
	case d.modif < 0:
		modif = fmt.Sprintf("%d", d.modif)
	}
	return count + "d" + sides + modSign + modif
}
