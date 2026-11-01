package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateFileName(t *testing.T) {
	type testCase struct {
		input string
		err   error
	}
	for _, tc := range []testCase{
		{"foo", assert.AnError},
		{"foo.sql", assert.AnError},
		{".sql", assert.AnError},
		{"", assert.AnError},
		{"foo.sql", assert.AnError},
		{"1234_foo.sql", nil},
	} {
		_, actual := validateFileName(tc.input)
		if tc.err != nil {
			assert.Error(t, actual)
		} else {
			assert.NoError(t, actual)
		}
	}
}
