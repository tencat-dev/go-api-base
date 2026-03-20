-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles
(
    id          UUID PRIMARY KEY      DEFAULT uuidv7(),
    name        VARCHAR(100) NOT NULL UNIQUE,

    description TEXT,
    is_system   BOOLEAN               DEFAULT FALSE,

    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE roles;
-- +goose StatementEnd
