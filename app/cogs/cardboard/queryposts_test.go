package cardboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tags_WithMagic_ShouldYieldRatingTags(t *testing.T) {
	q := NewQuery("xyz").WithMagic()

	q.WithSafe()
	actual := q.Tags()
	assert.Contains(t, actual, "xyz")
	assert.Contains(t, actual, noRatingS)
	assert.Contains(t, actual, noRatingE)

	q.WithUnsafe()
	actual = q.Tags()
	assert.Contains(t, actual, "xyz")
	// assert.Contains(t, actual, noRatingQ)
	assert.Contains(t, actual, noRatingG)
}

func Test_Tags_WithNoMagic_ShouldDropRatingTags(t *testing.T) {
	q := NewQuery("xyz").WithNoMagic()

	q.WithSafe()
	actual := q.Tags()
	assert.Contains(t, actual, "xyz")
	assert.NotContains(t, actual, noRatingQ)
	assert.NotContains(t, actual, noRatingE)
	assert.NotContains(t, actual, noRatingS)
	assert.NotContains(t, actual, noRatingG)

	q.WithUnsafe()
	actual = q.Tags()
	assert.Contains(t, actual, "xyz")
	assert.NotContains(t, actual, noRatingQ)
	assert.NotContains(t, actual, noRatingE)
	assert.NotContains(t, actual, noRatingS)
	assert.NotContains(t, actual, noRatingG)
}
