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
	testCases := []struct {
		path   string
		expect string
	}{
		{
			path:   "http://example.com/test?query=ignored",
			expect: "POST example.com/test",
		},
		{
			path:   "http://discord.com/api/v9/applications/964085462748774401/commands",
			expect: "POST discord.com/api/.+/applications/.+/commands",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.expect, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			actual := httpSpanNameFormatter("unit test", req)
			assert.Equal(t, tc.expect, actual)
		})
	}
}
