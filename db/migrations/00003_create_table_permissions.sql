-- +goose Up
-- +goose StatementBegin
CREATE TABLE permissions
(
    id     UUID PRIMARY KEY DEFAULT uuidv7(),
    object VARCHAR(100) NOT NULL,
    action VARCHAR(100) NOT NULL,
    UNIQUE (object, action)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE permissions;
-- +goose StatementEnd
