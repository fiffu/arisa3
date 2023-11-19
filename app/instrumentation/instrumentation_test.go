package instrumentation

import (
	"context"
	"net/http"
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

func Test_httpSpanNameFormatter(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com/test?query=ignored", nil)
	if err != nil {
		t.Fatal(err)
	}

	expect := "HTTP POST example.com/test"
	actual := httpSpanNameFormatter("ignoredOperation", req)
	assert.Equal(t, expect, actual)
}
