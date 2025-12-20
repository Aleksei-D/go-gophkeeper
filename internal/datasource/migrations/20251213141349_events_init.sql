-- +goose Up

CREATE TABLE IF NOT EXISTS events (
    login text NOT NULL,
    sync_date timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS events;
-- +goose StatementBegin
-- +goose StatementEnd