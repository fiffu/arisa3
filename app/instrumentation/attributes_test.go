package instrumentation

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPRequestPath_shouldOmitQuery(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://example.com/test?foo=bar", nil)
	if err != nil {
		t.Fatal(err)
	}

	attr := KV.HTTPRequestPath(req.Method, req.URL)
	assert.Equal(t, "POST http//example.com/test", attr.Value.AsString())
}
