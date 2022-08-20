package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/fiffu/arisa3/lib/functional"
	"github.com/rs/zerolog/log"
)

const (
	WildcardCharacter = "*"
	apiHost           = "danbooru.donmai.us"
	apiHostHTTPS      = "https://" + apiHost
	faviconPath       = "/favicon.ico"
)

var (
	MediaFileExts = []string{"png", "jpg", "jpeg", "gif"}
)

type client struct {
	auth        bool
	username    string
	apiKey      string
	host        string
	timeoutSecs int
	fetch       func(context.Context, *requests.Builder) error
}

func NewClient(username, apiKey string, timeoutSecs int) IClient {
	return newClient(username, apiKey, timeoutSecs)
}

func newClient(username, apiKey string, timeoutSecs int) *client {
	auth := false
	switch {
	case username != "" && apiKey != "":
		auth = true
	case username != "" && apiKey == "":
		log.Warn().Msg("Client username was provided, but API key was not. Skipping auth...")
	case username == "" && apiKey != "":
		log.Warn().Msg("Client API key was provided, but username was not. Skipping auth...")
	}

	if timeoutSecs == 0 {
		timeoutSecs = 2
	}
	return &client{
		auth,
		username,
		apiKey,
		apiHost,
		timeoutSecs,
		defaultFetcher,
	}
}

func commaJoin(strs []string) string {
	return strings.Join(strs, ",")
}

func spaceJoin(strs []string) string {
	return strings.Join(strs, " ")
}

func defaultFetcher(ctx context.Context, builder *requests.Builder) error {
	startTime := time.Now()
	var interceptor requests.RoundTripFunc = func(req *http.Request) (*http.Response, error) {
		reqID := newRequestID()
		logRequest(reqID, req, startTime)
		res, err := http.DefaultTransport.RoundTrip(req)
		logResponse(reqID, req, res, startTime, err)
		return res, err
	}
	return builder.Transport(interceptor).Fetch(ctx)
}

func newRequestID() string {
	seed := time.Now().UnixMicro()
	s := strconv.FormatInt(seed, 36)
	return strings.ToUpper(s)
}

func logRequest(reqID string, req *http.Request, startTime time.Time) {
	log.Info().
		Str("reqID", reqID).
		Msgf("%s %s", req.Method, req.URL.String())
}

func logResponse(reqID string, req *http.Request, res *http.Response, startTime time.Time, err error) {
	elapsed := time.Since(startTime)

	body, ioErr := io.ReadAll(res.Body)
	res.Body = io.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		log.Error().
			Str("reqID", reqID).
			Msgf("%s %s in %dms - request error: %s", req.Method, res.Status, elapsed.Milliseconds(), err)
	} else if ioErr != nil {
		log.Error().
			Str("reqID", reqID).
			Msgf("%s %s in %dms - io error: %s", req.Method, res.Status, elapsed.Milliseconds(), ioErr)
	} else {
		peek300Bytes := functional.SliceOf(body).Take(100) // First 100 bytes
		log.Info().
			Str("reqID", reqID).
			Msgf("%s %s in %dms - body: %s ...", req.Method, res.Status, elapsed.Milliseconds(), peek300Bytes)
	}
}

// UseAuth implements IClient.
func (c *client) UseAuth() bool {
	return c.auth
}

// FaviconURL implements IClient.
func (c *client) FaviconURL() string {
	return apiHostHTTPS + faviconPath
}

func (c *client) httpContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.timeoutSecs)*time.Second)
	return ctx, cancel
}

func (c *client) baseRequest() *requests.Builder {
	b := &requests.Builder{}
	b.Host(c.host)
	if c.UseAuth() {
		b.BasicAuth(c.username, c.apiKey)
	}
	return b
}

func (c *client) requestPosts() *requests.Builder {
	return c.baseRequest().Path("posts.json")
}

func (c *client) requestTags() *requests.Builder {
	return c.baseRequest().Path("tags.json")
}
