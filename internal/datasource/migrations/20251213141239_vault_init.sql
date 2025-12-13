-- +goose Up
CREATE TYPE data_type AS ENUM ('SECRET', 'TEXT', 'BLOB', 'CARD');
CREATE TABLE IF NOT EXISTS vault (
    login text NOT NULL,
    name text NOT NULL,
    metadata text,
    payload bytea NOT NULL,
    data_type data_type NOT NULL,
    comment text,
    update_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    is_deleted bool DEFAULT false,
    UNIQUE (login, name, data_type)
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS secrets;
-- +goose StatementBegin
-- +goose StatementEnd