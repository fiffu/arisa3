package cardboard

import (
	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/rs/zerolog/log"
)

func (d *domain) boringSearch(q IQueryPosts) ([]*api.Post, error) {
	return d.client.GetPosts(q.Tags())
}

func (d *domain) magicSearch(q IQueryPosts, tryGuessTerm bool) ([]*api.Post, error) {
	// Resolve any alias matches
	term, err := d.findAlias(q)
	if err != nil {
		// Log the error but don't break the flow
		log.Error().Err(err).Msgf("Errored while fetching aliases")
	} else {
		q.SetTerm(term)
	}

	// Perform search
	posts, err := d.boringSearch(q)
	if err != nil {
		return nil, err
	}

	// Filter results (reorder, omit, etc)
	filtered, err := d.filter(q, posts)
	if err != nil {
		// Log the error but don't break the flow
		log.Error().Err(err).Msgf("Errored while filtering results")
	} else {
		posts = filtered
	}

	switch {
	// Found results, we are done
	case len(posts) > 0:
		return posts, nil

	// If we still have a chance to guess, do it
	case tryGuessTerm:
		guess, err := d.guessTag(q)
		if err != nil {
			// Log the error then give up
			log.Error().Err(err).Msgf("Errored while filtering results")
			return posts, nil
		}
		if guess == q.Term() {
			// Our best guess exactly matched the query, give up
			return posts, nil
		} else {
			// Retry with guess
			q.SetTerm(guess)
			return d.magicSearch(q, false)
		}

	// No more guessing, give up
	default:
		return posts, nil
	}
}

func (d *domain) filter(q IQueryPosts, posts []*api.Post) ([]*api.Post, error) {
	opsMapping, err := d.repo.GetTagOperations()
	if err != nil {
		return nil, err
	}

	helper := &opsHelper{opsMapping}
	filters := []Filter{
		HasMediaFile(),
		HasURL(),
		Shuffle(),

		OmitFilter(helper),
		PromoteFilter(helper),
		DemoteFilter(helper),
	}
	return applyFilters(posts, filters), nil
}

func applyFilters(posts []*api.Post, filters []Filter) []*api.Post {
	for _, filter := range filters {
		posts = filter(posts)
		if len(posts) == 0 {
			break
		}
	}
	return posts
}

func (d *domain) findAlias(q IQueryPosts) (string, error) {
	term := q.Term()

	aliases, err := d.repo.GetAliases()
	if err != nil {
		return term, err
	}

	if actual, ok := aliases[Alias(term)]; ok {
		term = string(actual)
	}

	return term, nil
}

func (d *domain) guessTag(q IQueryPosts) (string, error) {
	term := q.Term()

	matches, err := d.client.GetTagsMatching(term + api.WildcardCharacter)
	if err != nil {
		return term, err
	}

	if len(matches) == 0 {
		return term, nil
	}

	return matches[0].Name, nil
}
