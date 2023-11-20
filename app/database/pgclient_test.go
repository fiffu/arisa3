package database

import (
	"context"
	"testing"

	"github.com/fiffu/arisa3/app/instrumentation"
	"github.com/stretchr/testify/assert"
)

func Test_newOperation(t *testing.T) {
	callable := func(ctx context.Context, query string, args ...any) (any, error) { return 1, nil }

	op := newOperation(callable, "test-caller")
	sql := "SELECT count(*) FROM my_table LIMIT 1"

	span := instrumentation.CaptureInstrumentation(t, func() {
		res, err := op(context.Background(), sql)
		assert.NoError(t, err)
		assert.Equal(t, 1, res)
	})

	assert.Equal(t, "DB: "+sql, span.Name)
	assert.Equal(t, "SELECT", span.Attributes.GetAsString("db_operation"))
	assert.Equal(t, sql, span.Attributes.GetAsString("db_query"))
}

func Test_firstWord(t *testing.T) {
	assert.Equal(t, "SELECT", firstWord("select count(*) from my_table limit 1"))
	assert.Equal(t, "COMMIT", firstWord("commit;"))
	assert.Equal(t, "", "")
}
