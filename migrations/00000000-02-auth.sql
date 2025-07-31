-- users
CREATE TABLE users (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);

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

-- providers
CREATE TABLE providers (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    jti TEXT,
    expiration TIMESTAMPTZ
);

CREATE UNIQUE INDEX providers_unique_idx ON providers (name);
CREATE INDEX providers_user_id_idx ON providers (user_id);

ALTER TABLE providers ADD CONSTRAINT fk_providers_user_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;

CREATE TRIGGER update_providers_updated BEFORE UPDATE ON providers
FOR EACH ROW EXECUTE FUNCTION update_updated_column();
