package cardboard

// queryPosts implements IQueryPosts
type queryPosts struct {
	magic   bool
	term    string
	rating  string
	guildID string
}

func NewQuery(term string) *queryPosts {
	return &queryPosts{
		term:   term,
		magic:  false,
		rating: "",
	}
}

// Methods for IQuery.

func (q *queryPosts) Tags() []string {
	tags := []string{q.term}
	if !q.MagicMode() && q.rating != "" {
		tags = append(tags, q.rating)
	}
	return tags
}

func (q *queryPosts) MagicMode() bool  { return q.magic }
func (q *queryPosts) Term() string     { return q.term }
func (q *queryPosts) GuildID() string  { return q.guildID }
func (q *queryPosts) SetTerm(s string) { q.term = s }

// Method-chaining builder.

func (q *queryPosts) SetNoMagic() *queryPosts           { q.magic = false; return q }
func (q *queryPosts) SetSafe() *queryPosts              { q.rating = tagRatingSafe; return q }
func (q *queryPosts) SetUnsafe() *queryPosts            { q.rating = tagRatingUnsafe; return q }
func (q *queryPosts) SetGuildID(gid string) *queryPosts { q.guildID = gid; return q }
