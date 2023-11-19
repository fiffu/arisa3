package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CaptureLogging(t *testing.T) {
	msg := "hello world"

	captured := CaptureLogging(t, func() {
		Infof(context.Background(), msg)
	})

	assert.Contains(t, captured, msg)
}
