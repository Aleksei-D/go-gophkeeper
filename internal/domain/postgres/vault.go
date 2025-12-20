package postgres

import (
	"context"
	"database/sql"
	"errors"
	"go-gophkeeper/internal/models"
	errors2 "go-gophkeeper/internal/utils/errors"
	"time"
)

const addDataQRY = `INSERT INTO vault as c1 
    (login, name, metadata, payload, data_type, comment, update_at, is_deleted) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (login, name, data_type) 
	DO UPDATE SET metadata = $3, payload = $4, comment = $6, update_at = $7, is_deleted = $8 
	WHERE EXCLUDED.update_at > c1.update_at`

type VaultRepository struct {
	db *sql.DB
}

func NewVaultRepository(db *sql.DB) *VaultRepository {
	return &VaultRepository{db: db}
}

func (c *VaultRepository) IsExist(ctx context.Context, name, owner string) (bool, error) {
	row := c.db.QueryRowContext(ctx, "SELECT * FROM vault WHERE login = $1 AND name = $2", owner, name)
	if err := row.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *VaultRepository) Get(ctx context.Context, name, login, dataType string) (models.VaultObject, error) {
	var vaultObject models.VaultObject
	var payload []byte
	var metadata string
	var comment string
	var updateAt time.Time

	query := `SELECT metadata, payload, comment, update_at FROM vault 
              WHERE name = $1 AND login = $2 AND data_type = $3  and is_deleted is false`
	err := c.db.QueryRowContext(
		ctx, query, name, login, dataType,
	).Scan(&metadata, &payload, &comment, &updateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return vaultObject, errors2.ErrNoContent
		}
		return vaultObject, err
	}

	vaultObject.Name = name
	vaultObject.Payload = payload
	vaultObject.Comment = comment
	vaultObject.UpdateAt = updateAt
	return vaultObject, err
}

func (c *VaultRepository) Add(ctx context.Context, vaultObject models.VaultObject) error {
	row := c.db.QueryRowContext(
		ctx,
		addDataQRY,
		vaultObject.Login,
		vaultObject.Name,
		vaultObject.Metadata,
		vaultObject.Payload,
		vaultObject.DataType,
		vaultObject.Comment,
		vaultObject.UpdateAt,
		vaultObject.IsDeleted,
	)
	return row.Err()
}

func (c *VaultRepository) GetDataToSync(ctx context.Context, login string) (models.VaultObjects, error) {
	var vaultObjects models.VaultObjects
	query := `SELECT name, metadata, payload, data_type, comment, update_at, is_deleted 
	FROM vault WHERE login = $1 AND update_at > (SELECT last_sync FROM events WHERE login = $1)`

	rows, err := c.db.QueryContext(ctx, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return vaultObjects, errors2.ErrNoContent
		}
		return vaultObjects, err
	}

	if err := rows.Err(); err != nil {
		return vaultObjects, err
	}
	defer rows.Close()
	for rows.Next() {
		var vaultObject models.VaultObject
		var name string
		var payload []byte
		var metadata string
		var dataType string
		var comment string
		var updateAt time.Time
		var isDeleted bool

		err := rows.Scan(&name, &metadata, &payload, &dataType, &comment, &comment, &updateAt, &isDeleted)
		if err != nil {
			return vaultObjects, err
		}
		vaultObject.Name = name
		vaultObject.Metadata = metadata
		vaultObject.Payload = payload
		vaultObject.DataType = dataType
		vaultObject.Comment = comment
		vaultObject.UpdateAt = updateAt
		vaultObject.IsDeleted = isDeleted
		vaultObjects = append(vaultObjects, vaultObject)
	}

	return vaultObjects, nil
}

func (c *VaultRepository) AddList(ctx context.Context, vaultObjects models.VaultObjects) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, vaultOBJ := range vaultObjects {
		row := c.db.QueryRowContext(
			ctx,
			addDataQRY,
			vaultOBJ.Login,
			vaultOBJ.Name,
			vaultOBJ.Metadata,
			vaultOBJ.Payload,
			vaultOBJ.DataType,
			vaultOBJ.Comment,
			vaultOBJ.UpdateAt,
			vaultOBJ.IsDeleted,
		)
		if err := row.Err(); err != nil {
			return err
		}
	}
	return tx.Commit()
}
