CREATE TABLE "colours" (
    userid TEXT,
    tstamp TIMESTAMP NOT NULL,  -- TEXT in database.TimeFormat (UTC, fixed width)
    reason TEXT  -- 'mutate' or 'reroll' or 'frozen'
);
