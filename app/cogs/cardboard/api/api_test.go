package api

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/fiffu/arisa3/testfixtures"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient_Auth(t *testing.T) {
	testCases := []struct {
		desc     string
		username string
		apiKey   string

		expectAuth bool
	}{
		{
			desc:       "no user, no apiKey should result in auth=false",
			username:   "",
			apiKey:     "",
			expectAuth: false,
		},
		{
			desc:       "no apiKey should result in auth=false",
			username:   "username",
			apiKey:     "",
			expectAuth: false,
		},
		{
			desc:       "no username should result in auth=false",
			username:   "",
			apiKey:     "apiKey",
			expectAuth: false,
		},
		{
			desc:       "username && password should result in auth=true",
			username:   "username",
			apiKey:     "apiKey",
			expectAuth: true,
		},
	}
	for i, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			timeout := i
			actual := NewClient(tc.username, tc.apiKey, timeout)

			assert.Equal(t, tc.expectAuth, actual.UseAuth())
		})
	}
}

func Test_maintenanceValidator(t *testing.T) {
	testCases := []struct {
		desc      string
		body      string
		expectErr error
	}{
		{
			desc:      "no maintenance",
			body:      "<html></html>",
			expectErr: nil,
		},
		{
			desc:      "have maintenance - h1",
			body:      `<html><h1>Danbooru is down for maintenance.</h1></html>`,
			expectErr: ErrUnderMaintenance,
		},
		{
			desc:      "have maintenance - title",
			body:      `<html><title>Downbooru</title></html>`,
			expectErr: ErrUnderMaintenance,
		},
		{
			desc:      "have maintenance - title",
			body:      testfixtures.Downbooru,
			expectErr: ErrUnderMaintenance,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			r := &http.Response{Body: io.NopCloser(strings.NewReader(tc.body))}
			err := (&client{}).maintenanceValidator(r)
			assert.Equal(t, tc.expectErr, err)
		})
	}
}

func Test_FaviconURL(t *testing.T) {
	client := NewClient("", "", 0)

	expect := apiHostHTTPS + faviconPath
	actual := client.FaviconURL()

	assert.Equal(t, expect, actual)
}
