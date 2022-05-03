package cardboard

import (
	"testing"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/stretchr/testify/assert"
)

func newTestOpsHelper(promote, demote, omit []string) *opsHelper {
	mapping := make(map[string]TagOperation)
	for _, tag := range promote {
		mapping[tag] = Promote
	}
	for _, tag := range demote {
		mapping[tag] = Demote
	}
	for _, tag := range omit {
		mapping[tag] = Omit
	}
	return &opsHelper{mapping}
}

func newTestPostWithTags(tags string) *api.Post {
	return &api.Post{AllTags: tags}
}

func assertEqualPostsByTag(t *testing.T, expect, actual []*api.Post) {
	t.Logf("Comparing length, expect: %d actual: %d", len(expect), len(actual))
	assert.Len(t, actual, len(expect))
	for i, e := range expect {
		a := actual[i]
		t.Logf("Comparing %v and %v", e, a)
		assert.Equal(t, *e, *a)
	}
}

func Test_hasOperation(t *testing.T) {
	h := newTestOpsHelper(
		[]string{"promote"},
		[]string{"demote"},
		[]string{},
	)
	type testCase struct {
		postTags                                     string
		expectPromoted, expectDemoted, expectOmitted bool
	}
	for _, tc := range []testCase{
		{"spam", false, false, false},
		{"promote", true, false, false},
		{"promote demote", true, true, false},
		{"", false, false, false},
	} {
		post := newTestPostWithTags(tc.postTags)
		assert.Equal(t, tc.expectPromoted, h.hasOperation(post, Promote))
		assert.Equal(t, tc.expectDemoted, h.hasOperation(post, Demote))
		assert.Equal(t, tc.expectOmitted, h.hasOperation(post, Omit))
	}
}

func Test_gatherByOperation(t *testing.T) {
	h := newTestOpsHelper(
		[]string{"promote"},
		[]string{"demote"},
		[]string{},
	)

	promotePost := newTestPostWithTags("promote")
	demotePost := newTestPostWithTags("demote")
	spam := newTestPostWithTags("spam")

	type testCase struct {
		groupBy       TagOperation
		expectGrouped []*api.Post
		expectRemain  []*api.Post
	}
	for _, tc := range []testCase{
		{
			Promote,
			[]*api.Post{promotePost},
			[]*api.Post{spam, spam, spam, demotePost},
		},
		{
			Demote,
			[]*api.Post{demotePost},
			[]*api.Post{promotePost, spam, spam, spam},
		},
		{
			Omit,
			[]*api.Post{},
			[]*api.Post{promotePost, spam, spam, spam, demotePost},
		},
	} {
		allPosts := []*api.Post{
			promotePost, spam, spam, spam, demotePost,
		}
		t.Log("Group by", tc.groupBy)
		group, remain := h.gatherByOperation(allPosts, tc.groupBy)
		assertEqualPostsByTag(t, tc.expectGrouped, group)
		assertEqualPostsByTag(t, tc.expectRemain, remain)

		group, remain = h.gatherByOperation(make([]*api.Post, 0), tc.groupBy)
		assert.Len(t, group, 0)
		assert.Len(t, remain, 0)
	}
}
