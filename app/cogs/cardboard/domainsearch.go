package cardboard

import (
	"context"
	"strings"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/fiffu/arisa3/app/log"
)

func (d *domain) boringSearch(ctx context.Context, q IQueryPosts) ([]*api.Post, error) {
	tags := q.Tags()
	log.Infof(ctx, "Querying for posts tagged='%v'", strings.Join(tags, " "))
	posts, err := d.client.GetPosts(ctx, tags)
	log.Infof(ctx, "Got %d posts, err=%v", len(posts), err)
	return posts, err
}

func (d *domain) magicSearch(ctx context.Context, q IQueryPosts, trySuggestion bool) ([]*api.Post, error) {
	log.Infof(ctx, "magicSearch for %+v", q)
	d.resolveAndNormalize(ctx, q)

	// Perform search
	posts, err := d.boringSearch(ctx, q)
	if err != nil {
		return nil, err
	}

	// Filter results (reorder, omit, etc)
	filtered, err := d.filter(ctx, q, posts)
	if err != nil {
		// Log the error but don't break the flow
		log.Errorf(ctx, err, "Errored while filtering results")
	} else {
		posts = filtered
	}

	// Found results, we are done
	if len(posts) > 0 {
		log.Infof(ctx, "magicSearch returning with %d posts", len(posts))
		return posts, nil
	}

	// If no results, yet we're not given a chance for suggestions, give up
	if !trySuggestion {
		return posts, nil
	}

	log.Infof(ctx, "magicSearch attempting to suggest another tag from query: %+v", q)
	suggest := d.fetchSuggestion(ctx, q)
	if suggest == "" {
		return posts, nil
	}

	log.Infof(ctx, "magicSearch retrying with suggestion=%+v", suggest)
	q.SetTerm(suggest)
	return d.magicSearch(ctx, q, false)
}

// resolveAndNormalize resolves alias matches and normalizes the query term into a tag.
func (d *domain) resolveAndNormalize(ctx context.Context, q IQueryPosts) {
	// Resolve any alias matches on the term.
	// Log any errors, but don't break the flow.
	if newTerm, err := d.findAlias(q); err != nil {
		log.Errorf(ctx, err, "Errored while fetching aliases")
	} else if newTerm != q.Term() {
		log.Infof(ctx, "Resolved alias %s -> %s", q.Term(), newTerm)
		q.SetTerm(newTerm)
	}

	// Convert spaces in the term into underscores
	taggified := taggify(q.Term())
	if taggified != q.Term() {
		log.Infof(ctx, "Taggify %s -> %s", q.Term(), taggified)
		q.SetTerm(taggified)
	}
	return
}

func (d *domain) fetchSuggestion(ctx context.Context, q IQueryPosts) string {
	guess, err := d.guessTag(ctx, q)
	if err != nil {
		// Log the error, but don't break the flow
		log.Errorf(ctx, err, "No suggestion: errored while fetching guesses")
		return ""
	}
	if guess == nil {
		log.Infof(ctx, "No suggestion: guess == nil")
		return ""
	}
	if guess.Name == q.Term() {
		log.Infof(ctx, "No suggestion: best guess already == %s", q.Term)
		return ""
	}
	return guess.Name
}

func (d *domain) filter(ctx context.Context, q IQueryPosts, posts []*api.Post) ([]*api.Post, error) {
	opsMapping := make(map[string]TagOperation)

	guildID := q.GuildID()
	if guildID != "" {
		guildOpsMapping, err := d.repo.GetTagOperations(guildID)
		if err != nil {
			return nil, err
		}
		opsMapping = guildOpsMapping
	}

	helper := &opsHelper{opsMapping}
	for fName, f := range map[string]Filter{
		"HasMediaFile": HasMediaFile(),
		"HasURL":       HasURL(),
		"Shuffle":      Shuffle(),

		"OmitFilter":    OmitFilter(helper),
		"PromoteFilter": PromoteFilter(helper),
		"DemoteFilter":  DemoteFilter(helper),
	} {
		before := len(posts)
		posts = f(posts)
		after := len(posts)

		if after < before {
			log.Infof(ctx, "%d posts excluded by filter %s", before-after, fName)
		}
		if after == 0 {
			break
		}
	}
	return posts, nil
}

func (d *domain) findAlias(q IQueryPosts) (string, error) {
	term := q.Term()
	aliases := make(map[Alias]Actual)

	guildID := q.GuildID()
	if guildID != "" {
		guildAliases, err := d.repo.GetAliases(guildID)
		if err != nil {
			return term, err
		}
		aliases = guildAliases
	}

	if actual, ok := aliases[Alias(term)]; ok {
		return string(actual), nil
	}
	return term, nil
}

func (d *domain) guessTag(ctx context.Context, q IQueryPosts) (*api.TagSuggestion, error) {
	term := q.Term()

	matches, err := d.client.AutocompleteTag(ctx, term)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, nil
	}
	match := matches[0]

	log.Infof(ctx, "Suggesting %#v from term: %s", match, term)
	return match, nil
}
