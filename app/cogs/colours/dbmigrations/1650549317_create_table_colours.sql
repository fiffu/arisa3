CREATE TABLE "colours" (
    userid TEXT,
    tstamp TIMESTAMP NOT NULL,
    reason TEXT,  -- 'mutate' or 'reroll' or 'frozen'
);
