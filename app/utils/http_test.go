package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewInstrumentedTransport_requestErrorDoesNotPanic(t *testing.T) {
	// A closed server refuses the connection, so there is no response.
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	client := &http.Client{Transport: NewInstrumentedTransport()}
	assert.NotPanics(t, func() {
		res, err := client.Get(url)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func Test_NewInstrumentedTransport_keepsBodyReadable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	client := &http.Client{Transport: NewInstrumentedTransport()}
	res, err := client.Get(srv.URL)
	assert.NoError(t, err)
	body, err := ReadAndReplaceBody(res)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(body))
}
