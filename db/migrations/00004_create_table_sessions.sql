-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions
(
    id            UUID PRIMARY KEY     DEFAULT uuidv7(),
    user_id       UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,

    refresh_jti   UUID        NOT NULL UNIQUE,
    token_version INT         NOT NULL DEFAULT 0,

    device_name   TEXT,
    user_agent    TEXT,
    ip_address    INET,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ NOT NULL,
    revoked_at    TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd