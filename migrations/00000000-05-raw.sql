-- gps history
CREATE TABLE raw (
    event_id BIGINT PRIMARY KEY REFERENCES events (id) ON DELETE CASCADE,
    data JSONB NOT NULL
);
