-- +goose Up
CREATE TABLE IF NOT EXISTS creds (
    user_login text NOT NULL,
    resource text NOT NULL,
    login text NOT NULL,
    password text NOT NULL
    update_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    comment text
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS creds;
-- +goose StatementBegin
-- +goose StatementEnd