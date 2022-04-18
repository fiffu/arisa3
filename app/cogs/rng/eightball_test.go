package rng

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_randChoice(t *testing.T) {
	choices := []string{"a", "b"}
	outcome := make(map[string]bool)

	i := 0
	for len(outcome) < 2 && i < 1000 {
		choice := randChoice(choices)
		outcome[choice] = true
	}
	assert.Len(t, outcome, 2)
}
