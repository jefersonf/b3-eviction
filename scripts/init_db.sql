
CREATE TABLE IF NOT EXISTS votes_minutely (
    bucket_minute TIMESTAMP NOT NULL, -- e.g. 2025-12-17 10:05:00
    eviction_id   TEXT NOT NULL,
    nominee_id    TEXT NOT NULL,
    votes         INT DEFAULT 0,
    PRIMARY KEY (bucket_minute, eviction_id, nominee_id)
);

-- Convert your existing table to a hypertable (partitioned by time)
SELECT create_hypertable('votes_minutely', 'bucket_minute');

CREATE MATERIALIZED VIEW votes_hourly
WITH (timescaledb.continuous) AS
SELECT
    time_bucket('1 hour', bucket_minute) as bucket_hour,
    eviction_id,
    nominee_id,
    SUM(votes) as total_votes
FROM votes_minutely
GROUP BY time_bucket('1 hour', bucket_minute), eviction_id, nominee_id;

ALTER MATERIALIZED VIEW votes_hourly SET (timescaledb.materialized_only = false);

-- Automatically refresh the view every minute
SELECT add_continuous_aggregate_policy('votes_hourly',
    start_offset => INTERVAL '3 days', -- Look back 3 days for changes
    end_offset => INTERVAL '1 minute', -- Keep 1 minute buffer for live data
    schedule_interval => INTERVAL '1 minute');




-- CREATE INDEX idx_votes_minutely_eviction ON votes_minutely (eviction_id);
-- CREATE INDEX idx_votes_minutely_nominee ON votes_minutely (nominee_id);
-- CREATE INDEX idx_votes_minutely_bucket ON votes_minutely (bucket_minute);


-- Case 1: Total Votes per Hour-Nominee
-- Super fast, scans very few rows
-- SELECT bucket_hour, nominee_id, total_votes 
-- FROM votes_hourly
-- ORDER BY bucket_hour DESC;

-- Case 2: Total Votes per Hour (All Nominees)
-- SELECT bucket_hour, SUM(total_votes) as hourly_total
-- FROM votes_hourly
-- GROUP BY bucket_hour
-- ORDER BY bucket_hour DESC;

-- Case 3: Grand Total Votes (The "Big Number")
-- This sums ~24 rows per day instead of ~1440 rows per day
-- SELECT SUM(total_votes) FROM votes_hourly;