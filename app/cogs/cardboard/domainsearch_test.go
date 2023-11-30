package cardboard

import (
	"context"
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
	mockClient.EXPECT().
		GetPosts(gomock.Any(), []string{"foo"}).
		Return(stubPosts, nil).
		Times(1)

	actualPosts, err := d.magicSearch(context.Background(), q, false)
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
	stubTags := []*api.TagSuggestion{{Name: "footlocker"}, {Name: "food"}}
	mockClient.EXPECT().GetPosts(gomock.Any(), []string{"foo"}).Return(noPosts, nil).Times(1)
	mockClient.EXPECT().AutocompleteTag(gomock.Any(), "foo").Return(stubTags, nil)
	mockClient.EXPECT().GetPosts(gomock.Any(), []string{"footlocker"}).Return(stubPosts, nil).Times(1)

	actualPosts, err := d.magicSearch(context.Background(), q, true)
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
	stubTags := []*api.TagSuggestion{{Name: "foo"}, {Name: "foot"}}
	mockClient.EXPECT().GetPosts(gomock.Any(), []string{"foo"}).Return(noPosts, nil).Times(1)
	mockClient.EXPECT().AutocompleteTag(gomock.Any(), "foo").Return(stubTags, nil)

	actualPosts, err := d.magicSearch(context.Background(), q, true)
	assert.NoError(t, err)
	assert.ElementsMatch(t, noPosts, actualPosts)
}

func Test_magicSearch_noSuggestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockIRepository(ctrl)
	mockClient := api.NewMockIClient(ctrl)
	d := &domain{mockRepo, mockClient}

	q := NewQuery("foo")
	noPosts := []*api.Post{}
	mockClient.EXPECT().GetPosts(gomock.Any(), []string{"foo"}).Return(noPosts, nil).Times(1)
	mockClient.EXPECT().AutocompleteTag(gomock.Any(), "foo").Return(nil, nil)

	actualPosts, err := d.magicSearch(context.Background(), q, true)
	assert.NoError(t, err)
	assert.ElementsMatch(t, noPosts, actualPosts)
}

func Test_magicSearch_guessErrorShouldBeSilent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := NewMockIRepository(ctrl)
	mockClient := api.NewMockIClient(ctrl)
	d := &domain{mockRepo, mockClient}

	q := NewQuery("foo")
	noPosts := []*api.Post{}
	var mockErr error = assert.AnError
	mockClient.EXPECT().GetPosts(gomock.Any(), []string{"foo"}).Return(noPosts, nil).Times(1)
	mockClient.EXPECT().AutocompleteTag(gomock.Any(), "foo").Return(nil, mockErr)

	actualPosts, err := d.magicSearch(context.Background(), q, true)
	assert.Nil(t, err, err)
	assert.Empty(t, actualPosts)
}
