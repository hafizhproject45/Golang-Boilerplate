-- Users
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL      PRIMARY KEY,
    name            VARCHAR,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ    NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);
