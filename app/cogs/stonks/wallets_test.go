package stonks

import (
	"context"
	"encoding/json"
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
	balances map[[2]string]int64  // {guild, user} -> balance
	active   map[[2]string]*React // {message, reactor} -> reaction, if active
	seen     map[[2]string]*React
	isBot    map[string]bool
}

func newFakeWalletRepo() *fakeWalletRepo {
	return &fakeWalletRepo{
		balances: map[[2]string]int64{},
		active:   map[[2]string]*React{},
		seen:     map[[2]string]*React{},
		isBot:    map[string]bool{},
	}
}

func (f *fakeWalletRepo) GuildEmojis(context.Context) (map[string]string, error) {
	return map[string]string{}, nil
}
func (f *fakeWalletRepo) SaveUser(_ context.Context, userID string, isBot bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.isBot[userID] = isBot
	return nil
}

func (f *fakeWalletRepo) AddCoinReaction(_ context.Context, r React) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	isBot, known := f.isBot[r.RecipientID]
	if !known {
		return false, ErrUnknownUser
	}
	if isBot {
		return false, nil
	}
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

// fakeUsers answers User() from a map, and counts the calls.
type fakeUsers struct {
	mu    sync.Mutex
	users map[string]*dgo.User
	err   error
	calls int
}

func (f *fakeUsers) User(userID string, _ ...dgo.RequestOption) (*dgo.User, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if u, ok := f.users[userID]; ok {
		return u, nil
	}
	return nil, &dgo.RESTError{Message: &dgo.APIErrorMessage{Code: dgo.ErrCodeUnknownUser, Message: "Unknown User"}}
}

const (
	botAuthor = "401"
	webhook   = "402"
)

func newTestWallets(now time.Time) (*wallets, *fakeWalletRepo, *fakeUsers) {
	repo := newFakeWalletRepo()
	emojis := &guildEmojis{}
	emojis.Set(guild, coin)
	users := &fakeUsers{users: map[string]*dgo.User{
		author:    {ID: author},
		reactor:   {ID: reactor},
		botAuthor: {ID: botAuthor, Bot: true},
	}}
	return &wallets{
		repo:   repo,
		emojis: emojis,
		now:    func() time.Time { return now },
	}, repo, users
}

func reaction(messageID, emojiID, userID string) *dgo.MessageReaction {
	return &dgo.MessageReaction{
		GuildID: guild, ChannelID: channel, MessageID: messageID, UserID: userID,
		Emoji: dgo.Emoji{ID: emojiID, Name: "birbcoin"},
	}
}

func reactionAdd(messageID, emojiID, userID, authorID string) *messageReactionAdd {
	return &messageReactionAdd{MessageReaction: *reaction(messageID, emojiID, userID), MessageAuthorID: authorID}
}

func Test_wallets_reactionRules(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	recent := snowflakeAt(now.Add(-time.Hour))
	edge := snowflakeAt(now.Add(-maxReactAge))
	old := snowflakeAt(now.Add(-maxReactAge - time.Second))

	cases := []struct {
		name      string
		r         *messageReactionAdd
		recipient string
		granted   bool
	}{
		{"coin emoji on a recent message", reactionAdd(recent, coin, reactor, author), author, true},
		{"message exactly one week old", reactionAdd(edge, coin, reactor, author), author, true},
		{"message older than one week", reactionAdd(old, coin, reactor, author), author, false},
		{"other emoji", reactionAdd(recent, "999", reactor, author), author, false},
		{"built-in emoji has no ID", reactionAdd(recent, "", reactor, author), author, false},
		{"self-react", reactionAdd(recent, coin, author, author), author, false},
		{"message from a bot", reactionAdd(recent, coin, reactor, botAuthor), botAuthor, false},
		{"message from a webhook", reactionAdd(recent, coin, reactor, webhook), webhook, false},
		{"payload without an author", reactionAdd(recent, coin, reactor, ""), "", false},
		{"guild without setup", func() *messageReactionAdd {
			r := reactionAdd(recent, coin, reactor, author)
			r.GuildID = "101"
			return r
		}(), author, false},
		{"DM", func() *messageReactionAdd {
			r := reactionAdd(recent, coin, reactor, author)
			r.GuildID = ""
			return r
		}(), author, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, repo, users := newTestWallets(now)
			w.addReaction(context.Background(), users, tc.r)

			bal, ok := repo.balance(tc.r.GuildID, tc.recipient)
			if tc.granted {
				assert.Equal(t, StartingBalance+CoinPerReaction, bal)
			} else {
				assert.False(t, ok, "no wallet change, got balance %d", bal)
			}
		})
	}
}

func Test_wallets_authorLookupIsSaved(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	w, repo, users := newTestWallets(now)
	react := func(w *wallets, i int, who string) {
		msg := snowflakeAt(now.Add(-time.Duration(i+1) * time.Minute))
		w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, who))
	}

	for i, who := range []string{author, author, webhook, webhook, botAuthor, botAuthor} {
		react(w, i, who)
	}
	assert.Equal(t, 3, users.calls, "one lookup per author, including a webhook and a bot")
	assert.Equal(t, map[string]bool{author: false, webhook: true, botAuthor: true}, repo.isBot)

	// After a restart, the saved flags still apply.
	restarted := &wallets{repo: repo, emojis: w.emojis, now: w.now}
	react(restarted, 10, author)
	react(restarted, 11, botAuthor)
	assert.Equal(t, 3, users.calls)

	bal, _ := repo.balance(guild, author)
	assert.Equal(t, StartingBalance+3*CoinPerReaction, bal)
	_, ok := repo.balance(guild, botAuthor)
	assert.False(t, ok, "a bot gets no wallet")
}

func Test_wallets_failedLookupIsNotSaved(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	w, repo, users := newTestWallets(now)
	msg := snowflakeAt(now.Add(-time.Minute))

	users.err = errors.New("503 Service Unavailable")
	w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, author))
	_, ok := repo.balance(guild, author)
	assert.False(t, ok)
	assert.Empty(t, repo.isBot)

	users.err = nil
	w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, author))
	bal, _ := repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal, "the next react looks the author up again")
}

func Test_wallets_payloadWithoutAuthorAsksNobody(t *testing.T) {
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	w, repo, users := newTestWallets(now)

	w.addReaction(context.Background(), users, reactionAdd(snowflakeAt(now.Add(-time.Minute)), coin, reactor, ""))

	assert.Equal(t, 0, users.calls)
	assert.Empty(t, repo.isBot, "no user row for an empty ID")
}

func Test_wallets_reactUnreactReact(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	msg := snowflakeAt(now.Add(-time.Hour))
	w, repo, users := newTestWallets(now)

	w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, author))
	w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, author)) // a duplicate event
	bal, _ := repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal, "one pair gives one coin")

	w.removeReaction(ctx, reaction(msg, coin, reactor))
	w.removeReaction(ctx, reaction(msg, coin, reactor))
	bal, _ = repo.balance(guild, author)
	assert.Equal(t, StartingBalance, bal, "the un-react takes the coin back once")

	w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, author))
	bal, _ = repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal)
}

func Test_wallets_unreactOnOldMessageKeepsCoin(t *testing.T) {
	ctx := context.Background()
	day2 := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	msg := snowflakeAt(day2.Add(-24 * time.Hour))
	w, repo, users := newTestWallets(day2)

	w.addReaction(ctx, users, reactionAdd(msg, coin, reactor, author))
	w.now = func() time.Time { return day2.Add(7 * 24 * time.Hour) }
	w.removeReaction(ctx, reaction(msg, coin, reactor))

	bal, _ := repo.balance(guild, author)
	assert.Equal(t, StartingBalance+CoinPerReaction, bal)
}

// The payload uses the field names from Discord's docs for MESSAGE_REACTION_ADD.
// It is not a captured event.
func Test_messageReactionAdd_decodesAuthor(t *testing.T) {
	raw := `{"user_id":"500","channel_id":"200","message_id":"600","guild_id":"100",
		"emoji":{"id":"300","name":"birbcoin"},"message_author_id":"400","burst":false,"type":0}`
	var r messageReactionAdd
	require.NoError(t, json.Unmarshal([]byte(raw), &r))
	assert.Equal(t, "400", r.MessageAuthorID)
	assert.Equal(t, "500", r.UserID)
	assert.Equal(t, "100", r.GuildID)
	assert.Equal(t, "300", r.Emoji.ID)
}

func Test_snowflakeAt(t *testing.T) {
	at := time.Date(2026, 11, 10, 12, 0, 0, 0, time.UTC)
	got, err := dgo.SnowflakeTimestamp(snowflakeAt(at))
	require.NoError(t, err)
	assert.True(t, at.Equal(got))
}
