-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    src_link VARCHAR(2048) NOT NULL UNIQUE,
    short_link VARCHAR(128) NOT NULL UNIQUE,
    visits INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX src_short_index ON links(src_link, short_link);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE links;
-- +goose StatementEnd
