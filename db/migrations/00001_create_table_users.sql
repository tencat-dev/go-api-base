-- +goose Up
-- +goose StatementBegin
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id            UUID PRIMARY KEY     DEFAULT uuidv7(),

    name          TEXT        NOT NULL,
    email         TEXT        NOT NULL UNIQUE,

    password_hash TEXT        NOT NULL,
    token_version INT         NOT NULL DEFAULT 0,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
