package database

import (
	"context"
	"testing"

	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/stretchr/testify/assert"
)

func Test_newOperation(t *testing.T) {
	callable := func(ctx context.Context, query string, args ...any) (any, error) { return 1, nil }

	op := newOperation[any](callable, "test-caller", "test-operation")

	span := instrumentation.CaptureInstrumentation(t, func() {
		res, err := op(context.Background(), "SELECT 1")
		assert.NoError(t, err)
		assert.Equal(t, 1, res)
	})

	assert.Equal(t, "test-operation", span.Attributes.GetAsString("db_operation"))
	assert.Equal(t, "SELECT 1", span.Attributes.GetAsString("db_query"))

}
