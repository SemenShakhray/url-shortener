-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url (
    id SERIAL PRIMARY KEY,
    alias VARCHAR NOT NULL UNIQUE,
    url VARCHAR NOT NULL UNIQUE,
    create_at TIMESTAMP(0),
    update_at TIMESTAMP(0)
);
CREATE INDEX IF NOT EXISTS url_alias ON url(alias);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS url_alias;
DROP TABLE IF EXISTS url;

-- +goose StatementEnd
