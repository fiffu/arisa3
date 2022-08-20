package api

import (
	"testing"

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

func Test_FaviconURL(t *testing.T) {
	client := NewClient("", "", 0)

	expect := apiHostHTTPS + faviconPath
	actual := client.FaviconURL()

	assert.Equal(t, expect, actual)
}
