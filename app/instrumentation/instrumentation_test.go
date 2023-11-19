package instrumentation

import (
	"context"
	"testing"

	"github.com/fiffu/arisa3/app/log"
	"github.com/stretchr/testify/assert"
)

func Test_NewInstrumentationClient(t *testing.T) {

	var c Client
	var err error
	logs := log.CaptureLogging(t, func() {
		c, err = NewInstrumentationClient(context.Background())
	})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Contains(t, logs, "Instrumentation client enabled: false serviceName:")
}
