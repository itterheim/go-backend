-- gps history
CREATE TABLE locations_history (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION NOT NULL,
    provider_id BIGINT,
    user_id BIGINT NOT NULL
);

CREATE INDEX locations_history_timestamp_idx ON locations_history (timestamp);
CREATE INDEX locations_history_lat_idx ON locations_history (latitude);
CREATE INDEX locations_history_lon_idx ON locations_history (longitude);
CREATE INDEX locations_history_user_id_idx ON locations_history (user_id);
CREATE INDEX locations_history_provider_id_idx ON locations_history (provider_id);

-- locations
CREATE TABLE locations_places (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    note TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius DOUBLE PRECISION,
    user_id BIGINT NOT NULL,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX locations_places_name_unique_idx ON locations_places (name);
CREATE INDEX locations_places_lat_idx ON locations_places (latitude);
CREATE INDEX locations_places_lon_idx ON locations_places (longitude);
CREATE INDEX locations_places_user_id_idx ON locations_places (user_id);

CREATE TRIGGER update_locations_places_updated BEFORE UPDATE ON locations_places
FOR EACH ROW EXECUTE FUNCTION update_updated_column();

-- join table for M:N relation between history and places
CREATE TABLE locations_history_places (
    history_id BIGINT NOT NULL,
    place_id BIGINT NOT NULL,
    PRIMARY KEY (history_id, place_id)
);

CREATE INDEX locations_history_places_history_id_idx ON locations_history_places (history_id);
CREATE INDEX locations_history_places_place_id_idx ON locations_history_places (place_id);

ALTER TABLE locations_history_places ADD CONSTRAINT fk_locations_history_places_history_id FOREIGN KEY (history_id) REFERENCES locations_history (id) ON DELETE CASCADE;
ALTER TABLE locations_history_places ADD CONSTRAINT fk_locations_history_places_place_id FOREIGN KEY (place_id) REFERENCES locations_history (id) ON DELETE CASCADE;