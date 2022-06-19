package cardboard

import (
	"strings"

	"github.com/fiffu/arisa3/app/types"
)

func taggify(s string) string {
	from := " "
	into := "_"
	return strings.ReplaceAll(s, from, into)
}

func untaggify(s string) string {
	from := "_"
	into := " "
	return strings.ReplaceAll(s, from, into)
}

func getGuildID(req types.ICommandEvent) string {
	ixn := req.Interaction()
	if ixn == nil {
		return ""
	}
	return ixn.GuildID
}
