-- +goose Up
CREATE TABLE IF NOT EXISTS cards (
    login text NOT NULL,
    file_name text NOT NULL,
    update_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    is_udpated bool,
    comment text
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS cards;
-- +goose StatementBegin
-- +goose StatementEnd