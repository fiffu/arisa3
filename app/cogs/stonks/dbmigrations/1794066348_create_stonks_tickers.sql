-- One row per ticker. prices[i] is the price at slots[i].
-- prices: thousandths of a coin per unit. slots: Unix seconds at the start of each update slot.
CREATE TABLE "stonks_tickers" (
    ticker TEXT PRIMARY KEY,
    prices BIGINT[] NOT NULL DEFAULT '{}',
    slots  BIGINT[] NOT NULL DEFAULT '{}',
    meta   JSONB    NOT NULL DEFAULT '{}',
    CHECK (cardinality(prices) = cardinality(slots))
);
