-- users
CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE INDEX users_username_idx ON users (username);
CREATE UNIQUE INDEX users_username_unique_idx ON users (username);

CREATE TRIGGER update_users_updated BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_column();

-- tokens
CREATE TABLE tokens (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT,
    jti TEXT NOT NULL,
    expiration TIMESTAMPTZ NOT NULL,
    blocked BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX tokens_user_id_idx ON tokens (user_id);

ALTER TABLE tokens ADD CONSTRAINT fk_tokens_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE TRIGGER update_tokens_updated BEFORE UPDATE ON tokens
FOR EACH ROW EXECUTE FUNCTION update_updated_column();

-- devices
CREATE TABLE devices (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    name TEXT NOT NULL,
    description TEXT,
    jti TEXT,
    expiration TIMESTAMPTZ
);

CREATE TRIGGER update_devices_updated BEFORE UPDATE ON devices
FOR EACH ROW EXECUTE FUNCTION update_updated_column();
