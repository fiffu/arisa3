CREATE TABLE "colours_log" (
    userid   TEXT,
    username TEXT,
    colour   TEXT,
    reason   TEXT,  -- 'mutate' or 'reroll' or 'frozen'
    tstamp   TIMESTAMP NOT NULL,
) PARTITION BY RANGE (tstamp);

-- 2-year sliding window view of the table
CREATE VIEW "colours_logview" AS
    SELECT * FROM "colours_log"
    WHERE tstamp > current_timestamp - INTERVAL '2 years';
