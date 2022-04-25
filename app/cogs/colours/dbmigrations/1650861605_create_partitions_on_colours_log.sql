-- Ref:
-- https://www.postgresql.org/docs/14/ddl-partitioning.html

-- Let's have yearly partitions (no reason why but I don't plan to revisit DBs more than once a year)
CREATE TABLE "colours_log_2022" PARTITION OF colours_log FOR VALUES FROM ('2022-01-01') TO ('2023-01-01');
CREATE TABLE "colours_log_2023" PARTITION OF colours_log FOR VALUES FROM ('2023-01-01') TO ('2024-01-01');
CREATE TABLE "colours_log_2024" PARTITION OF colours_log FOR VALUES FROM ('2024-01-01') TO ('2025-01-01');
CREATE TABLE "colours_log_2025" PARTITION OF colours_log FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
CREATE TABLE "colours_log_2026" PARTITION OF colours_log FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');
CREATE TABLE "colours_log_2027" PARTITION OF colours_log FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');
CREATE TABLE "colours_log_2028" PARTITION OF colours_log FOR VALUES FROM ('2028-01-01') TO ('2029-01-01');
CREATE TABLE "colours_log_2029" PARTITION OF colours_log FOR VALUES FROM ('2029-01-01') TO ('2030-01-01');
CREATE TABLE "colours_log_2030" PARTITION OF colours_log FOR VALUES FROM ('2030-01-01') TO ('2031-01-01');

-- I don't really think we will get here...
CREATE TABLE "colours_log_future" PARTITION OF colours_log FOR VALUES FROM ('2031-01-01') TO ('2222-01-01');

-- Index the partition key, since the docs say "in most scenarios it is helpful"
CREATE INDEX ON colours_log (tstamp);
