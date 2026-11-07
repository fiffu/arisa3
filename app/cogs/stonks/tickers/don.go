package tickers

import (
	"context"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fiffu/arisa3/app/cogs/stonks/tickers/api"
	"golang.org/x/net/html"
)

// Don uses DefaultSlotSize when it has never read a window.
const DefaultSlotSize = 12 * time.Hour

type DonMeta struct {
	// CoveredUntil is the end of the last window that Don read. The next window starts
	// here, not at the previous slot. After a failed fetch, the update writes no patch,
	// so the next window also covers the failed one, and its posts still count.
	CoveredUntil *time.Time `json:"covered_until,omitempty"`
}

// Don uses the average text length of the posts since the previous update as the price.
// With no new posts, the price stays the same.
type Don struct {
	Client *api.TrumpsTruth
}

func (d *Don) Ticker() string { return "DON" }

func (d *Don) Fetch(ctx context.Context, st State[DonMeta]) (Result[DonMeta], error) {
	posts, err := d.Client.Posts(ctx)
	if err != nil {
		return Result[DonMeta]{}, err
	}

	from := st.PrevSlot
	if st.Meta.CoveredUntil != nil {
		from = *st.Meta.CoveredUntil
	}
	if from.IsZero() {
		from = st.Slot.Add(-DefaultSlotSize)
	}
	to := st.Slot
	patch := DonMeta{CoveredUntil: &to}

	total, count := 0, 0
	for _, post := range posts {
		if post.At.Before(from) || !post.At.Before(st.Slot) {
			continue
		}
		n := postLength(post.HTML)
		if n == 0 {
			continue // media-only post
		}
		total += n
		count++
	}

	if count == 0 {
		return Result[DonMeta]{Raw: st.PrevPrice, Outcome: NoNewData, Patch: patch}, nil
	}
	return Result[DonMeta]{Raw: float64(total) / float64(count), Outcome: OK, Patch: patch}, nil
}

// blockTags separate words, so "<p>a</p><p>b</p>" reads as "a b".
var blockTags = map[string]bool{
	"p": true, "br": true, "div": true, "blockquote": true, "li": true,
}

// The archive links each image and video to its attachment store, with link text such as
// "View image". Those links are not post text, so they do not count.
func postLength(desc string) int {
	z := html.NewTokenizer(strings.NewReader(desc))
	var text strings.Builder
	inMediaLink := false

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken: // end of input
			return utf8.RuneCountInString(strings.Join(strings.Fields(text.String()), " "))

		case html.StartTagToken, html.EndTagToken, html.SelfClosingTagToken:
			name, hasAttr := z.TagName()
			tag := string(name)
			switch {
			case tag == "a" && tt == html.StartTagToken:
				inMediaLink = hasAttr && isMediaLink(z)
			case tag == "a" && tt == html.EndTagToken:
				inMediaLink = false
			case blockTags[tag]:
				text.WriteString(" ")
			}

		case html.TextToken:
			// Text() returns the text with entities such as &amp; decoded.
			if !inMediaLink {
				text.Write(z.Text())
			}
		}
	}
}

// isMediaLink reads the attributes of the current <a> tag.
func isMediaLink(z *html.Tokenizer) bool {
	for {
		key, val, more := z.TagAttr()
		if string(key) == "href" {
			u, err := url.Parse(string(val))
			return err == nil && strings.Contains(u.Path, "/attachments/")
		}
		if !more {
			return false
		}
	}
}
