-- events
CREATE TABLE events (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    until TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL,
    tags TEXT[],
    note TEXT,
    reference TEXT,
    provider_id BIGINT,
    user_id BIGINT NOT NULL
);

CREATE INDEX events_type_idx ON events (type);
CREATE INDEX events_timestamp_idx ON events (timestamp);
CREATE INDEX events_until_idx ON events (until);
CREATE INDEX events_status_idx ON events (status);
CREATE INDEX events_tags_idx ON events USING GIN(tags);
CREATE INDEX events_reference_idx ON events (reference);
CREATE INDEX events_user_id_idx ON events (user_id);
CREATE INDEX events_provider_id_idx ON events (provider_id);

ALTER TABLE events ADD CONSTRAINT fk_events_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
ALTER TABLE events ADD CONSTRAINT fk_events_provider_id FOREIGN KEY (provider_id) REFERENCES providers (id) ON DELETE SET NULL;

-- tags
CREATE TABLE tags (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tag TEXT NOT NULL,
    description TEXT,
    parent_id BIGINT,
    user_id BIGINT NOT NULL
);

CREATE INDEX tags_tag_trgm ON tags USING GIN(tag gin_trgm_ops);
CREATE INDEX tags_parent_id ON tags (parent_id);
CREATE INDEX tags_user_id_idx ON tags (user_id);

ALTER TABLE tags ADD CONSTRAINT fk_tags_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;