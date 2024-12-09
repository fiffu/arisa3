package lib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FrozenClock(t *testing.T) {
	now := time.Now()
	fc := FrozenClock(t, now)

	assert.Equal(t, now, fc.Now())

	positive := 5 * time.Minute
	fc.Add(positive)
	assert.Equal(t, now.Add(positive), fc.Now())

	negative := -2 * time.Second
	fc.Add(negative)
	assert.Equal(t, now.Add(positive).Add(negative), fc.Now())

	now2 := time.Now().Add(30 * time.Millisecond)
	fc.Set(now2)
	assert.Equal(t, now2, fc.Now())
}
