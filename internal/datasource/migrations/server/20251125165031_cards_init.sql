-- +goose Up
CREATE TABLE IF NOT EXISTS cards (
    login text NOT NULL,
    card_number text NOT NULL,
    expirationm_month INT NOT NULL,
    expiration_year INT NOT NULL,
    update_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    comment text
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS cards;
-- +goose StatementBegin
-- +goose StatementEnd