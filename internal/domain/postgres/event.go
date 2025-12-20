package postgres

import (
	"context"
	"database/sql"
	"errors"
	errors2 "go-gophkeeper/internal/utils/errors"
	"time"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (e *EventRepository) Get(ctx context.Context, login string) (time.Time, error) {
	var date time.Time
	err := e.db.QueryRowContext(
		ctx, "SELECT sync_date FROM events WHERE login = $1", login,
	).Scan(&date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return date, errors2.ErrNoContent
		}
		return date, err
	}

	return date, err
}

func (e *EventRepository) Add(ctx context.Context, login string) error {
	insertQRY := `INSERT INTO events (login, sync_date) VALUES ($1, CURRENT_TIMESTAMP)
    ON CONFLICT (login) DO UPDATE SET sync_date = CURRENT_TIMESTAMP`

	row := e.db.QueryRowContext(ctx, insertQRY, login)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
