package cardboard

type Safety int

const (
	undefined Safety = 0
	safe      Safety = 1
	unsafe    Safety = 2
)

var (
	safeTags   = []string{"rating:g"}
	unsafeTags = []string{"-rating:g", "-rating:s"}
)

// queryPosts implements IQueryPosts
type queryPosts struct {
	magic   bool
	term    string
	safety  Safety
	guildID string
}

func NewQuery(term string) *queryPosts {
	return &queryPosts{
		term:   term,
		magic:  false,
		safety: undefined,
	}
}

// Methods for IQuery.

func (q *queryPosts) Tags() []string {
	tags := []string{q.term}
	if q.MagicMode() {
		switch q.safety {
		case safe:
			tags = append(tags, safeTags...)
		case unsafe:
			tags = append(tags, unsafeTags...)
		}
	}
	return tags
}

func (q *queryPosts) MagicMode() bool  { return q.magic }
func (q *queryPosts) Term() string     { return q.term }
func (q *queryPosts) GuildID() string  { return q.guildID }
func (q *queryPosts) SetTerm(s string) { q.term = s }

// Method-chaining builder.

func (q *queryPosts) WithMagic() *queryPosts             { q.magic = true; return q }
func (q *queryPosts) WithNoMagic() *queryPosts           { q.magic = false; return q }
func (q *queryPosts) WithSafe() *queryPosts              { q.safety = safe; return q }
func (q *queryPosts) WithUnsafe() *queryPosts            { q.safety = unsafe; return q }
func (q *queryPosts) WithGuildID(gid string) *queryPosts { q.guildID = gid; return q }
