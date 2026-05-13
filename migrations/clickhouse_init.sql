CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.test_events (
                                                     event_id UUID,
                                                     user_id UUID,
                                                     test_id UUID,
                                                     topic String,
                                                     score UInt32,
                                                     max_score UInt32,
                                                     percentage Float64,
                                                     correct UInt32,
                                                     incorrect UInt32,
                                                     duration_ms UInt64,
                                                     created_at DateTime
) ENGINE = MergeTree()
    ORDER BY (created_at, user_id, test_id);