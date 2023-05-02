package log

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type stringable struct{ x string }

func (t stringable) String() string { return t.x }

func Test_stringify(t *testing.T) {
	testCases := []struct {
		input  any
		expect string
	}{
		{
			input:  "plain string",
			expect: "plain string",
		},
		{
			input:  time.Unix(1683008000, 0),
			expect: "2023-05-02T06:13:20Z",
		},
		{
			input:  errors.New("some error"),
			expect: "some error",
		},
		{
			input:  stringable{"stringify"},
			expect: "stringify",
		},
		{
			input:  123,
			expect: "123",
		},
		{
			input:  123.4,
			expect: "123.4",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			assert.Equal(t, tc.expect, stringify(tc.input))
		})
	}
}
