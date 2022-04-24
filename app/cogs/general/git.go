package general

import (
	"fmt"

	"github.com/fiffu/arisa3/app/types"
)

func (c *Cog) gitCommand() *types.Command {
	webURL := c.cfg.RepoWebURL
	return types.NewCommand("git").ForChat().
		Desc("Spot a bug? Want to contribute? >> " + webURL).
		Handler(c.git)
}

func (c *Cog) git(req types.ICommandEvent) error {
	var (
		repoName  = c.cfg.RepoName
		webURL    = c.cfg.RepoWebURL
		issuesURL = c.cfg.RepoIssuesURL
		gitURL    = c.cfg.RepoGitCloneURL
	)

	issues := fmt.Sprintf("Bugs and suggestions: create an issue!\n%s", issuesURL)
	contribs := fmt.Sprintf("Contributions: pull requests welcome!\n```%s```", gitURL)
	embed := types.NewEmbed().
		Title(repoName).
		URL(webURL).
		Description(issues + "\n\n" + contribs)

	resp := types.NewResponse().Embeds(embed)
	return req.Respond(resp)
}
