package stonks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fiffu/arisa3/app/log"

	dgo "github.com/bwmarrin/discordgo"
)

// maxReactAge is the maximum message age for a react or un-react to count.
// The age comes from the message ID, so the check needs no API call.
const maxReactAge = 7 * 24 * time.Hour

// guildEmojis caches stonks_guilds.data.emoji_id, so a reaction event needs no query
// unless it uses the coin emoji. Only /stonks-setup changes it.
type guildEmojis struct {
	mu sync.RWMutex
	m  map[string]string
}

func (g *guildEmojis) Get(guildID string) (string, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	id, ok := g.m[guildID]
	return id, ok
}

func (g *guildEmojis) Set(guildID, emojiID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.m == nil {
		g.m = map[string]string{}
	}
	g.m[guildID] = emojiID
}

// userGetter is the part of *dgo.Session that wallets uses. Tests replace it.
type userGetter interface {
	User(userID string, options ...dgo.RequestOption) (*dgo.User, error)
}

// messageReactionAdd is the MESSAGE_REACTION_ADD payload. discordgo (v0.27 to v0.29, and master
// in Feb 2026) does not decode message_author_id, so wallets reads it from the raw event.
type messageReactionAdd struct {
	dgo.MessageReaction
	MessageAuthorID string `json:"message_author_id"`
}

type wallets struct {
	repo   IWalletRepository
	emojis *guildEmojis
	now    func() time.Time
}

func newReact(r *dgo.MessageReaction, msgAuthor string) React {
	return React{
		EmojiID:     r.Emoji.ID,
		GuildID:     r.GuildID,
		ChannelID:   r.ChannelID,
		MessageID:   r.MessageID,
		ReactorID:   r.UserID,
		Source:      SourceLive,
		RecipientID: msgAuthor,
	}
}

// age reads the send time from the message ID, so it needs no API call.
// An ID that does not parse gives an age above maxReactAge, so the react does not count.
func (r React) age(now time.Time) time.Duration {
	sent, err := dgo.SnowflakeTimestamp(r.MessageID)
	if err != nil {
		return maxReactAge + 1
	}
	return now.Sub(sent)
}

// onMessageReactionAdd takes the raw event, because only the raw payload has the message author.
func (w *wallets) onMessageReactionAdd(ctx context.Context, s *dgo.Session, e *dgo.Event) {
	var r messageReactionAdd
	if err := json.Unmarshal(e.RawData, &r); err != nil {
		log.Errorf(ctx, err, "stonks: failed to decode %s", e.Type)
		return
	}
	w.addReaction(ctx, s, &r)
}

func (w *wallets) addReaction(ctx context.Context, api userGetter, r *messageReactionAdd) {
	react := newReact(&r.MessageReaction, r.MessageAuthorID)
	coinEmojiID, configured := w.emojis.Get(react.GuildID)
	if !configured {
		return
	}

	isCoinEmoji := react.EmojiID == coinEmojiID
	isRecent := react.age(w.now()) <= maxReactAge
	// message_author_id is optional in the payload, and "" is not a valid user ID in SQL.
	hasAuthor := react.RecipientID != ""
	differentUser := react.RecipientID != react.ReactorID
	if !(isCoinEmoji && isRecent && hasAuthor && differentUser) {
		return
	}

	granted, err := w.repo.AddCoinReaction(ctx, react)
	if errors.Is(err, ErrUnknownUser) {
		// Ask Discord only for an author that no coin react has seen before.
		if err = w.lookupUser(ctx, api, react.RecipientID); err == nil {
			granted, err = w.repo.AddCoinReaction(ctx, react)
		}
	}
	if err != nil {
		log.Errorf(ctx, err, "stonks: failed to add coin reaction on message %s", r.MessageID)
		return
	}
	if granted {
		log.Infof(ctx, "stonks: %s gave a coin to %s for message %s", react.ReactorID, react.RecipientID, r.MessageID)
	}
}

func (w *wallets) lookupUser(ctx context.Context, api userGetter, recipientID string) error {
	isBot, err := lookUpIsBot(api, recipientID)
	if err != nil {
		return fmt.Errorf("lookup recipient %s: %w", recipientID, err)
	}
	if err := w.repo.SaveUser(ctx, recipientID, isBot); err != nil {
		return fmt.Errorf("saving recipient %s: %w", recipientID, err)
	}
	return nil
}

func lookUpIsBot(api userGetter, userID string) (bool, error) {
	u, err := api.User(userID)
	var restErr *dgo.RESTError
	switch {
	case errors.As(err, &restErr) && restErr.Message != nil && restErr.Message.Code == dgo.ErrCodeUnknownUser:
		// A webhook message has the webhook ID as its author, and that ID is not a user.
		return true, nil
	case err != nil:
		return false, err
	}
	return u.Bot, nil
}

func (w *wallets) onMessageReactionRemove(ctx context.Context, _ *dgo.Session, e *dgo.MessageReactionRemove) {
	w.removeReaction(ctx, e.MessageReaction)
}

func (w *wallets) removeReaction(ctx context.Context, r *dgo.MessageReaction) {
	// The un-react event has no message author. RemoveCoinReaction reads the recipient from the stored row.
	react := newReact(r, "")
	coinEmojiID, configured := w.emojis.Get(react.GuildID)
	if !configured {
		return
	}

	isCoinEmoji := react.EmojiID == coinEmojiID
	isRecent := react.age(w.now()) <= maxReactAge
	if !(isCoinEmoji && isRecent) {
		return
	}

	removed, err := w.repo.RemoveCoinReaction(ctx, react.MessageID, react.ReactorID)
	if err != nil {
		log.Errorf(ctx, err, "stonks: failed to remove coin reaction on message %s", react.MessageID)
		return
	}
	if removed {
		log.Infof(ctx, "stonks: %s took back a coin for message %s", react.ReactorID, react.MessageID)
	}
}
