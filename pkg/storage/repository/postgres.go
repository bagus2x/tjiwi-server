package repository

import (
	"context"
	"database/sql"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/bagus2x/tjiwi/pkg/storage"
	"github.com/sirupsen/logrus"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) storage.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, st *model.Storage) error {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			INSERT INTO
				Storage
				(supervisor_id, name, description, is_deleted, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6)
			RETURNING
				id
	`
	err := tx.QueryRowContext(
		ctx,
		query,
		st.Supervisor.ID,
		st.Name,
		st.Description,
		st.IsDeleted,
		st.CreatedAt,
		st.UpdatedAt,
	).Scan(&st.ID)

	return err
}

func (r *repository) FindByID(ctx context.Context, storageID int64) (*model.Storage, error) {
	query := `
			SELECT
				s.id, p.id, p.photo, p.username, p.email, s.name, s.description, s.created_at,
				s.updated_at
			FROM
				Storage s
			JOIN
				Profile p
			ON
				s.supervisor_id = p.id	
			WHERE
				s.id = $1 AND s.is_deleted = FALSE
	`

	var s model.Storage

	err := r.db.QueryRowContext(ctx, query, storageID).Scan(
		&s.ID,
		&s.Supervisor.ID,
		&s.Supervisor.Photo,
		&s.Supervisor.Username,
		&s.Supervisor.Email,
		&s.Name,
		&s.Description,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewError(nil, app.ENotFound)
		}

		return nil, err
	}

	return &s, nil
}

func (r *repository) FindBySupervisorID(ctx context.Context, supervisorID int64) ([]*model.Storage, error) {
	query := `
			SELECT
				id, supervisor_id, name, description, is_deleted, created_at, updated_at
			FROM
				Storage
			WHERE
				supervisor_id = $1 AND is_deleted = FALSE
			ORDER BY
				updated_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, supervisorID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	storages := make([]*model.Storage, 0)

	for rows.Next() {
		var s model.Storage

		err := rows.Scan(
			&s.ID,
			&s.Supervisor.ID,
			&s.Name,
			&s.Description,
			&s.IsDeleted,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		storages = append(storages, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return storages, nil
}

func (r *repository) Update(ctx context.Context, st *model.Storage) error {
	query := `
			UPDATE
				Storage
			SET
				name = $1,
				description	= $2,
				updated_at = $3
			WHERE
				id = $4
	`

	res, err := r.db.ExecContext(ctx, query, st.Name, st.Description, st.UpdatedAt, st.ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return app.NewError(nil, app.ENotFound)
	}

	return nil
}

func (r *repository) SoftDelete(ctx context.Context, storageID int64, isDeleted bool) error {
	query := `
			UPDATE
				Storage
			SET
				is_deleted = $1
			WHERE
				id = $2
	`

	res, err := r.db.ExecContext(ctx, query, isDeleted, storageID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return app.NewError(nil, app.ENotFound)
	}

	return nil
}

func (r *repository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	c := context.WithValue(ctx, db.TransactionKey{}, tx)
	err = fn(c)
	if err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			logrus.Error("Failed to rollback transaction", errTx)
		}
		return err
	}

	if errTX := tx.Commit(); errTX != nil {
		logrus.Error("Failed to commmit transaction", errTX)
	}

	return nil
}
