-- This was yearly Postgres partitions of colours_log. SQLite has no partitioning and the table is small,
-- so it is one table with an index on the timestamp.
CREATE INDEX "colours_log_tstamp" ON "colours_log" (tstamp);
CREATE INDEX "colours_log_userid" ON "colours_log" (userid, tstamp);
