CREATE TABLE "colours_log" (
    userid   TEXT,
    username TEXT,
    colour   TEXT,
    reason   TEXT,  -- 'mutate' or 'reroll' or 'frozen'
    tstamp   TIMESTAMP NOT NULL  -- TEXT in database.TimeFormat (UTC, fixed width)
);

-- 2-year sliding window view of the table.
-- strftime matches database.TimeFormat's layout, so the string comparison is a time comparison.
CREATE VIEW "colours_logview" AS
    SELECT * FROM "colours_log"
    WHERE tstamp > strftime('%Y-%m-%dT%H:%M:%fZ', 'now', '-2 years');
