-- data holds the guild settings and the backfill progress, for example:
--   {"emoji_id": "123", "market_channel_id": "456", "backfill": {"<channel_id>": {...}}}
-- Change single keys with || or jsonb_set, never by rewriting the whole object.
CREATE TABLE "stonks_guilds" (
    guild_id   BIGINT      PRIMARY KEY,
    data       JSONB       NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- balance is in thousandths of a unit. The row with ticker 'coin' holds the coin balance.
-- 'coin' is not a row in stonks_tickers, so ticker has no foreign key.
CREATE TABLE "stonks_wallets" (
    guild_id   BIGINT      NOT NULL,
    user_id    BIGINT      NOT NULL,
    ticker     TEXT        NOT NULL,
    balance    BIGINT      NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (guild_id, user_id, ticker)
);

-- One row per (message, reactor). active changes on each react and un-react,
-- so a pair grants at most one coin at any time.
CREATE TABLE "stonks_coin_reactions" (
    guild_id     BIGINT      NOT NULL,
    channel_id   BIGINT      NOT NULL,
    message_id   BIGINT      NOT NULL,
    reactor_id   BIGINT      NOT NULL,
    recipient_id BIGINT      NOT NULL,
    active       BOOLEAN     NOT NULL,
    source       TEXT        NOT NULL CHECK (source IN ('live', 'backfill')),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (message_id, reactor_id)
);

-- For top donors and top recipients in /coins.
CREATE INDEX stonks_coin_reactions_recipient ON stonks_coin_reactions (guild_id, recipient_id) WHERE active;
CREATE INDEX stonks_coin_reactions_reactor   ON stonks_coin_reactions (guild_id, reactor_id)   WHERE active;
