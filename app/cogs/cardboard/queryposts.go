package cardboard

import "strings"

const (
	safe   = "rating:s"
	unsafe = "-rating:s"
)

// queryPosts implements IQueryPosts
type queryPosts struct {
	magic  bool
	term   string
	rating string
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
func (q *queryPosts) SetTerm(s string) { q.term = s }
func (q *queryPosts) String() string   { return strings.Join(q.Tags(), " ") }

// Method-chaining builder.

func (q *queryPosts) NoMagic() *queryPosts   { q.magic = false; return q }
func (q *queryPosts) SetSafe() *queryPosts   { q.rating = safe; return q }
func (q *queryPosts) SetUnsafe() *queryPosts { q.rating = unsafe; return q }
