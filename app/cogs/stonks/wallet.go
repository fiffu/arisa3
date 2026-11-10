package stonks

import (
	"context"
	"sync"
	"time"

	"github.com/fiffu/arisa3/app/log"

	dgo "github.com/bwmarrin/discordgo"
)

// ReactWindow is the maximum message age for a react or un-react to count.
// The age comes from the message ID, so the check needs no API call.
const ReactWindow = 7 * 24 * time.Hour

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

// fetchMessage is dgo.Session.ChannelMessage. Tests replace it.
type fetchMessage func(channelID, messageID string) (*dgo.Message, error)

type wallets struct {
	repo   IWalletRepository
	emojis *guildEmojis
	fetch  fetchMessage
	now    func() time.Time
}

// counts checks the rules that need no API call: a guild with a coin emoji,
// that emoji, and a message inside ReactWindow.
func (w *wallets) counts(r *dgo.MessageReaction) bool {
	if r.GuildID == "" || r.Emoji.ID == "" {
		return false
	}
	if emojiID, ok := w.emojis.Get(r.GuildID); !ok || emojiID != r.Emoji.ID {
		return false
	}
	sent, err := dgo.SnowflakeTimestamp(r.MessageID)
	if err != nil {
		return false
	}
	return w.now().Sub(sent) <= ReactWindow
}

func (w *wallets) onReactionAdd(ctx context.Context, r *dgo.MessageReaction) {
	if !w.counts(r) {
		return
	}
	// The reaction event in this discordgo version has no message author, so fetch the message.
	msg, err := w.fetch(r.ChannelID, r.MessageID)
	if err != nil {
		log.Errorf(ctx, err, "stonks: failed to fetch message %s for a coin reaction", r.MessageID)
		return
	}
	if msg.Author == nil || msg.Author.Bot || msg.WebhookID != "" || msg.Author.ID == r.UserID {
		return
	}

	granted, err := w.repo.AddCoinReaction(ctx, CoinReaction{
		GuildID:     r.GuildID,
		ChannelID:   r.ChannelID,
		MessageID:   r.MessageID,
		ReactorID:   r.UserID,
		RecipientID: msg.Author.ID,
		Source:      SourceLive,
	})
	if err != nil {
		log.Errorf(ctx, err, "stonks: failed to add coin reaction on message %s", r.MessageID)
		return
	}
	if granted {
		log.Infof(ctx, "stonks: %s gave a coin to %s for message %s", r.UserID, msg.Author.ID, r.MessageID)
	}
}

func (w *wallets) onReactionRemove(ctx context.Context, r *dgo.MessageReaction) {
	if !w.counts(r) {
		return
	}
	revoked, err := w.repo.RemoveCoinReaction(ctx, r.MessageID, r.UserID)
	if err != nil {
		log.Errorf(ctx, err, "stonks: failed to remove coin reaction on message %s", r.MessageID)
		return
	}
	if revoked {
		log.Infof(ctx, "stonks: %s took back a coin for message %s", r.UserID, r.MessageID)
	}
}
