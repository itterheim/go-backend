-- gps history
CREATE TABLE locations_history (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    accuracy DOUBLE PRECISION NOT NULL,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX locations_history_lat_idx ON locations_history (latitude);
CREATE INDEX locations_history_lon_idx ON locations_history (longitude);

-- locations
CREATE TABLE locations_places (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT NOT NULL,
    tags TEXT[],
    note TEXT,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    radius DOUBLE PRECISION,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX locations_places_name_unique_idx ON locations_places (name);
CREATE INDEX locations_places_lat_idx ON locations_places (latitude);
CREATE INDEX locations_places_lon_idx ON locations_places (longitude);
CREATE INDEX locations_places_tags_idx ON locations_places USING GIN(tags);

CREATE TRIGGER update_locations_places_updated BEFORE UPDATE ON locations_places
FOR EACH ROW EXECUTE FUNCTION update_updated_column();