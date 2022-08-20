package cardboard

import (
	"testing"

	"github.com/fiffu/arisa3/app/cogs/cardboard/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_magicSearch_tryGuessTermIsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockIRepository(ctrl)
	mockClient := api.NewMockIClient(ctrl)
	d := &domain{mockRepo, mockClient}

	q := NewQuery("foo")
	stubPosts := []*api.Post{{File: "https://foo.jpg", FileExt: "jpg"}}
	mockClient.EXPECT().GetPosts([]string{"foo"}).Return(stubPosts, nil).Times(1)

	actualPosts, err := d.magicSearch(q, false)
	assert.NoError(t, err)
	assert.ElementsMatch(t, stubPosts, actualPosts)
}

func Test_magicSearch_tryGuessTermWithRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockIRepository(ctrl)
	mockClient := api.NewMockIClient(ctrl)
	d := &domain{mockRepo, mockClient}

	q := NewQuery("foo")
	noPosts := []*api.Post{}
	stubPosts := []*api.Post{{File: "https://food.jpg", FileExt: "jpg"}}
	stubTags := []*api.Tag{{Name: "footlocker"}, {Name: "food"}}
	mockClient.EXPECT().GetPosts([]string{"foo"}).Return(noPosts, nil).Times(1)
	mockClient.EXPECT().GetPosts([]string{"food"}).Return(stubPosts, nil).Times(1)
	mockClient.EXPECT().GetTagsMatching("foo*").Return(stubTags, nil)

	actualPosts, err := d.magicSearch(q, true)
	assert.NoError(t, err)
	assert.ElementsMatch(t, stubPosts, actualPosts)
}

func Test_magicSearch_tryGuessTermNoRetry(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockIRepository(ctrl)
	mockClient := api.NewMockIClient(ctrl)
	d := &domain{mockRepo, mockClient}

	q := NewQuery("foo")
	noPosts := []*api.Post{}
	stubTags := []*api.Tag{{Name: "foo"}, {Name: "foot"}}
	mockClient.EXPECT().GetPosts([]string{"foo"}).Return(noPosts, nil).Times(1)
	mockClient.EXPECT().GetTagsMatching("foo*").Return(stubTags, nil)

	actualPosts, err := d.magicSearch(q, true)
	assert.NoError(t, err)
	assert.ElementsMatch(t, noPosts, actualPosts)
}

func Test_guessTag(t *testing.T) {
	testCases := []struct {
		desc          string
		term          string
		stubResponse  []string
		stubError     error
		expectGuessed string
		expectError   error
	}{
		{
			desc:          "No matches",
			term:          "ham",
			stubResponse:  []string{},
			expectGuessed: "ham",
		},
		{
			desc:          "Exact match",
			term:          "ham",
			stubResponse:  []string{"hammer", "ham", "hamstring"},
			expectGuessed: "ham",
		},
		{
			desc:          "Shortest candidate",
			term:          "ha",
			stubResponse:  []string{"hamper", "ham", "hamstring"},
			expectGuessed: "ham",
		},
		{
			desc:        "Error",
			term:        "ham",
			stubError:   assert.AnError,
			expectError: assert.AnError,
		},
	}

	ctrl := gomock.NewController(t)
	mockClient := api.NewMockIClient(ctrl)
	d := &domain{nil, mockClient}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			stubbed := make([]*api.Tag, 0)
			for _, tagName := range tc.stubResponse {
				stubbed = append(stubbed, &api.Tag{Name: tagName})
			}
			mockClient.EXPECT().GetTagsMatching(tc.term+api.WildcardCharacter).
				Return(stubbed, tc.stubError).
				Times(1)

			q := NewQuery(tc.term)
			actual, err := d.guessTag(q)

			if tc.expectError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectGuessed, actual)
			}
		})
	}
}
