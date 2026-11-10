package stonks

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeWalletRepo keeps coin balances in memory, with the same rules as the SQL.
type fakeWalletRepo struct {
	mu       sync.Mutex
	settings map[string][2]string        // guild -> {emoji, channel}
	balances map[[2]string]int64         // {guild, user} -> balance
	active   map[[2]string]*CoinReaction // {message, reactor} -> reaction, if active
	seen     map[[2]string]*CoinReaction
}

func newFakeWalletRepo() *fakeWalletRepo {
	return &fakeWalletRepo{
		settings: map[string][2]string{},
		balances: map[[2]string]int64{},
		active:   map[[2]string]*CoinReaction{},
		seen:     map[[2]string]*CoinReaction{},
	}
}

func (f *fakeWalletRepo) GuildEmojis(context.Context) (map[string]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := map[string]string{}
	for g, s := range f.settings {
		out[g] = s[0]
	}
	return out, nil
}

func (f *fakeWalletRepo) SaveGuildSettings(_ context.Context, guildID, emojiID, channelID string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.settings[guildID] = [2]string{emojiID, channelID}
	return nil
}

func (f *fakeWalletRepo) AddCoinReaction(_ context.Context, r CoinReaction) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	wallet := [2]string{r.GuildID, r.RecipientID}
	if _, ok := f.balances[wallet]; !ok {
		f.balances[wallet] = StartingBalance
	}
	key := [2]string{r.MessageID, r.ReactorID}
	if f.active[key] != nil {
		return false, nil
	}
	if prev := f.seen[key]; prev != nil {
		r = *prev // a re-react keeps the stored recipient
	}
	f.active[key], f.seen[key] = &r, &r
	f.balances[[2]string{r.GuildID, r.RecipientID}] += CoinPerReaction
	return true, nil
}

func (f *fakeWalletRepo) RemoveCoinReaction(_ context.Context, messageID, reactorID string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := [2]string{messageID, reactorID}
	r := f.active[key]
	if r == nil {
		return false, nil
	}
	delete(f.active, key)
	f.balances[[2]string{r.GuildID, r.RecipientID}] -= CoinPerReaction
	return true, nil
}

func (f *fakeWalletRepo) balance(guild, user string) (int64, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	b, ok := f.balances[[2]string{guild, user}]
	return b, ok
}

// snowflakeAt returns a message ID that Discord would give a message sent at t.
func snowflakeAt(t time.Time) string {
	return strconv.FormatInt((t.UnixMilli()-1420070400000)<<22, 10)
}

const (
	guild   = "100"
	channel = "200"
	coin    = "300"
	author  = "400"
	reactor = "500"
)

func newTestWallets(now time.Time, msgs map[string]*dgo.Message) (*wallets, *fakeWalletRepo) {
	repo := newFakeWalletRepo()
	emojis := &guildEmojis{}
	emojis.Set(guild, coin)
	return &wallets{
		repo:   repo,
		emojis: emojis,
		now:    func() time.Time { return now },
		fetch: func(_, messageID string) (*dgo.Message, error) {
			if m, ok := msgs[messageID]; ok {
				return m, nil
			}
			return nil, errors.New("404 Not Found")
		},
	}, repo
}

func reaction(messageID, emojiID, userID string) *dgo.MessageReaction {
	return &dgo.MessageReaction{
		GuildID: guild, ChannelID: channel, MessageID: messageID, UserID: userID,
		Emoji: dgo.Emoji{ID: emojiID, Name: "birbcoin"},
	}
}

func Test_wallets_reactionRules(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	recent := snowflakeAt(now.Add(-time.Hour))
	edge := snowflakeAt(now.Add(-ReactWindow))
	old := snowflakeAt(now.Add(-ReactWindow - time.Second))
	byBot := snowflakeAt(now.Add(-2 * time.Hour))
	byWebhook := snowflakeAt(now.Add(-3 * time.Hour))
	missing := snowflakeAt(now.Add(-4 * time.Hour))
	msgs := map[string]*dgo.Message{
		recent:    {ID: recent, Author: &dgo.User{ID: author}},
		edge:      {ID: edge, Author: &dgo.User{ID: author}},
		old:       {ID: old, Author: &dgo.User{ID: author}},
		byBot:     {ID: byBot, Author: &dgo.User{ID: author, Bot: true}},
		byWebhook: {ID: byWebhook, Author: &dgo.User{ID: author}, WebhookID: "9"},
	}

	cases := []struct {
		name    string
		r       *dgo.MessageReaction
		granted bool
	}{
		{"coin emoji on a recent message", reaction(recent, coin, reactor), true},
		{"message exactly one week old", reaction(edge, coin, reactor), true},
		{"message older than one week", reaction(old, coin, reactor), false},
		{"other emoji", reaction(recent, "999", reactor), false},
		{"built-in emoji has no ID", reaction(recent, "", reactor), false},
		{"self-react", reaction(recent, coin, author), false},
		{"message from a bot", reaction(byBot, coin, reactor), false},
		{"message from a webhook", reaction(byWebhook, coin, reactor), false},
		{"message not found", reaction(missing, coin, reactor), false},
		{"guild without setup", func() *dgo.MessageReaction {
			r := reaction(recent, coin, reactor)
			r.GuildID = "101"
			return r
		}(), false},
		{"DM", func() *dgo.MessageReaction {
			r := reaction(recent, coin, reactor)
			r.GuildID = ""
			return r
		}(), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, repo := newTestWallets(now, msgs)
			w.onReactionAdd(context.Background(), tc.r)

			bal, ok := repo.balance(tc.r.GuildID, author)
			if tc.granted {
				assert.Equal(t, StartingBalance+CoinPerReaction, bal)
			} else {
				assert.False(t, ok, "no wallet change, got balance %d", bal)
			}
		})
	}
}

func Test_wallets_reactUnreactReact(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	msg := snowflakeAt(now.Add(-time.Hour))
	w, repo := newTestWallets(now, map[string]*dgo.Message{msg: {ID: msg, Author: &dgo.User{ID: author}}})

	w.onReactionAdd(ctx, reaction(msg, coin, reactor))
	w.onReactionAdd(ctx, reaction(msg, coin, reactor)) // a duplicate event
	bal, _ := repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal, "one pair gives one coin")

	w.onReactionRemove(ctx, reaction(msg, coin, reactor))
	w.onReactionRemove(ctx, reaction(msg, coin, reactor))
	bal, _ = repo.balance(guild, author)
	assert.Equal(t, StartingBalance, bal, "the un-react takes the coin back once")

	w.onReactionAdd(ctx, reaction(msg, coin, reactor))
	bal, _ = repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal)
}

func Test_wallets_unreactOnOldMessageKeepsCoin(t *testing.T) {
	ctx := context.Background()
	day2 := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	msg := snowflakeAt(day2.Add(-24 * time.Hour))
	w, repo := newTestWallets(day2, map[string]*dgo.Message{msg: {ID: msg, Author: &dgo.User{ID: author}}})

	w.onReactionAdd(ctx, reaction(msg, coin, reactor))
	w.now = func() time.Time { return day2.Add(7 * 24 * time.Hour) }
	w.onReactionRemove(ctx, reaction(msg, coin, reactor))

	bal, _ := repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal)
}

func Test_parseEmojiID(t *testing.T) {
	cases := []struct {
		in, id string
		ok     bool
	}{
		{"<:birbcoin:1234567890>", "1234567890", true},
		{" <a:birbspin:42> ", "42", true},
		{"🪙", "", false},
		{":birbcoin:", "", false},
		{"<:birbcoin:12> extra", "", false},
	}
	for _, tc := range cases {
		id, ok := parseEmojiID(tc.in)
		assert.Equal(t, tc.ok, ok, tc.in)
		assert.Equal(t, tc.id, id, tc.in)
	}
}

func Test_snowflakeAt(t *testing.T) {
	at := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	got, err := dgo.SnowflakeTimestamp(snowflakeAt(at))
	require.NoError(t, err)
	assert.True(t, at.Equal(got))
}
