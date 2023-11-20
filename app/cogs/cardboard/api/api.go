package api

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/carlmjohnson/requests"
	"github.com/fiffu/arisa3/app/log"
	"github.com/fiffu/arisa3/app/utils"
)

const (
	WildcardCharacter = "*"
	apiHost           = "danbooru.donmai.us"
	apiHostHTTPS      = "https://" + apiHost
	faviconPath       = "/favicon.ico"
)

var (
	MediaFileExts = []string{"png", "jpg", "jpeg", "gif"}

	ErrUnderMaintenance = errors.New("API is under maintenance")
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
	ctx := context.Background()

	auth := false
	switch {
	case username != "" && apiKey != "":
		auth = true
	case username != "" && apiKey == "":
		log.Infof(ctx, "Client username was provided, but API key was not. Skipping auth...")
	case username == "" && apiKey != "":
		log.Infof(ctx, "Client API key was provided, but username was not. Skipping auth...")
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
	return builder.
		Transport(utils.NewInstrumentedTransport()).
		Fetch(ctx)
}

// UseAuth implements IClient.
func (c *client) UseAuth() bool {
	return c.auth
}

// FaviconURL implements IClient.
func (c *client) FaviconURL() string {
	return apiHostHTTPS + faviconPath
}

func (c *client) httpContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(c.timeoutSecs)*time.Second)
	return ctx, cancel
}

func (c *client) baseRequest() *requests.Builder {
	b := &requests.Builder{}
	b.Host(c.host)
	if c.UseAuth() {
		b.BasicAuth(c.username, c.apiKey)
	}
	b.AddValidator(c.maintenanceValidator)
	return b
}

func (c *client) maintenanceValidator(r *http.Response) error {
	body, err := utils.ReadAndReplaceBody(r)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	maint, err := isMaintenanceDoc(doc)
	if err != nil {
		return err
	}

	if maint {
		return ErrUnderMaintenance
	}
	return nil
}

func isMaintenanceDoc(doc *goquery.Document) (bool, error) {
	elems := map[string]string{
		"downbooru":            doc.Find("title").Text(),
		"down for maintenance": doc.Find("h1").Text(),
	}
	indicators := 0.0
	for expect, elem := range elems {
		if strings.Contains(strings.ToLower(elem), expect) {
			indicators += 1
		}
	}
	probability := indicators / float64(len(elems))
	return probability >= 0.5, nil
}

func (c *client) postsResource() *requests.Builder {
	return c.baseRequest().Path("posts.json")
}

func (c *client) tagsResource() *requests.Builder {
	return c.baseRequest().Path("tags.json")
}

// Not a REST endpoint.
func (c *client) autocompleteEndpoint() *requests.Builder {
	return c.baseRequest().Path("autocomplete")
}
