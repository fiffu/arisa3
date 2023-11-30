package cardboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tags_WithMagic_ShouldYieldRatingTags(t *testing.T) {
	q := NewQuery("xyz").WithMagic()

	q.WithSafe()
	tags := []string{"xyz"}
	actual := q.Tags()
	assert.Equal(t, append(tags, safeTags...), actual)

	q.WithUnsafe()
	actual = q.Tags()
	assert.Equal(t, append(tags, unsafeTags...), actual)
}

func Test_Tags_WithNoMagic_ShouldDropRatingTags(t *testing.T) {
	q := NewQuery("xyz").WithNoMagic()

	q.WithSafe()
	actual := q.Tags()
	assert.Contains(t, actual, "xyz")
	assert.NotContains(t, actual, unsafeTags)
	assert.NotContains(t, actual, safeTags)

	q.WithUnsafe()
	actual = q.Tags()
	assert.Contains(t, actual, "xyz")
	assert.NotContains(t, actual, unsafeTags)
	assert.NotContains(t, actual, safeTags)
}
