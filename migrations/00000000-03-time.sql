-- events
CREATE TABLE events (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    until TIMESTAMPTZ,

    status VARCHAR(20) NOT NULL,
    tags TEXT[],
    note TEXT,

    provider_id BIGINT,
    user_id BIGINT NOT NULL
);

CREATE INDEX events_type_idx ON events (type);
CREATE INDEX events_timestamp_idx ON events (timestamp);
CREATE INDEX events_until_idx ON events (until);
CREATE INDEX events_status_idx ON events (status);
CREATE INDEX events_tags_idx ON events USING GIN(tags);
CREATE INDEX events_user_id_idx ON events (user_id);
CREATE INDEX events_provider_id_idx ON events (provider_id);

ALTER TABLE events ADD CONSTRAINT fk_events_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL;
ALTER TABLE events ADD CONSTRAINT fk_events_provider_id FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE SET NULL;

-- actions
CREATE TABLE actions (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    event_id BIGINT NOT NULL,
    reference VARCHAR(63),
    reference_id BIGINT,
    tags TEXT[],
    note TEXT
)

CREATE INDEX actions_event_id_idx ON actions (event_id);
CREATE INDEX actions_reference_idx ON actions (reference);
CREATE INDEX actions_reference_id_idx ON actions (reference_id);
CREATE INDEX actions_tags_idx ON actions USING GIN(tags);
