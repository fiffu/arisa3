package cardboard

import (
	"strings"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/rs/zerolog/log"
)

func (d *domain) boringSearch(q IQueryPosts) ([]*api.Post, error) {
	tags := q.Tags()
	log.Info().Msgf("Querying for posts tagged='%v'", strings.Join(tags, " "))
	posts, err := d.client.GetPosts(tags)
	log.Info().Msgf("Got %d posts, err=%v", len(posts), err)
	return posts, err
}

func (d *domain) magicSearch(q IQueryPosts, tryGuessTerm bool) ([]*api.Post, error) {
	log.Info().Msgf("magicSearch for %+v", q)

	// Resolve any alias matches on the term
	newTerm, err := d.findAlias(q)
	if err != nil {
		// Log the error but don't break the flow
		log.Error().Err(err).Msgf("Errored while fetching aliases")
	} else if newTerm != q.Term() {
		log.Info().Msgf("Resolved alias %s -> %s", q.Term(), newTerm)
		q.SetTerm(newTerm)
	}

	// Convert spaces in the term into underscores
	newTerm = taggify(q.Term())
	if newTerm != q.Term() {
		log.Info().Msgf("Taggify %s -> %s", q.Term(), newTerm)
		q.SetTerm(newTerm)
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
		log.Info().Msgf("magicSearch returning with %d posts", len(posts))
		return posts, nil

	// If we still have a chance to guess, do it
	case tryGuessTerm:
		log.Info().Msgf("magicSearch attempting to guess another tag from query: %+v", q)
		guess, err := d.guessTag(q)
		if err != nil {
			// Log the error then give up
			log.Error().Err(err).Msgf("Errored while best-guessing query")
			return posts, nil
		}
		if guess == q.Term() {
			// Our best guess exactly matched the query, give up
			return posts, nil
		} else {
			// Retry with guess
			log.Info().Msgf("magicSearch retrying with guess=%s", guess)
			q.SetTerm(guess)
			return d.magicSearch(q, false)
		}

	// No more guessing, give up
	default:
		return posts, nil
	}
}

func (d *domain) filter(q IQueryPosts, posts []*api.Post) ([]*api.Post, error) {
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
			log.Info().Msgf("%d posts excluded by filter %s", before-after, fName)
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

func (d *domain) guessTag(q IQueryPosts) (string, error) {
	term := q.Term()

	matches, err := d.client.GetTagsMatching(term + api.WildcardCharacter)
	if err != nil {
		return term, err
	}

	if len(matches) == 0 {
		return term, nil
	}

	matches = api.TagsSorter{Data: matches, Compare: api.ByTagLength}.Sorted()
	return matches[0].Name, nil
}
