package repository

import (
	"context"
	"database/sql"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/bagus2x/tjiwi/pkg/storagemember"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) storagemember.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, sm *model.StorageMember) error {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			INSERT INTO
				Storage_Member
				(storage_id, member_id, is_admin, is_active, is_deleted, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT
				(storage_id, member_id)
			DO UPDATE SET
				is_deleted = FALSE,
				is_admin = $3,
				is_active = $4
			RETURNING
				id
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		sm.Storage.ID,
		sm.Member.ID,
		sm.IsAdmin,
		sm.IsActive,
		sm.IsDeleted,
		sm.CreatedAt,
		sm.UpdatedAt,
	).Scan(&sm.ID)
	if err != nil {
		return err
	}

	return err
}

func (r *repository) FindByID(ctx context.Context, storMembID int64) (*model.StorageMember, error) {
	query := `
			SELECT
				id, storage_id, member_id, is_admin, is_active, is_deleted, created_at, updated_at
			FROM
				Storage_Member
			WHERE
				id = $1 AND is_deleted = FALSE
			FOR UPDATE
	`

	var sm model.StorageMember

	err := r.db.QueryRowContext(ctx, query, storMembID).Scan(
		&sm.ID,
		&sm.Storage.ID,
		&sm.Member.ID,
		&sm.IsAdmin,
		&sm.IsActive,
		&sm.IsDeleted,
		&sm.CreatedAt,
		&sm.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewError(nil, app.ENotFound)
		}

		return nil, err
	}

	return &sm, nil
}

func (r *repository) FindByStorageIDAndUserID(ctx context.Context, storageID, userID int64) (*model.StorageMember, error) {
	query := `
			SELECT
				id, storage_id, member_id, is_admin, is_active, is_deleted, created_at, updated_at
			FROM
				Storage_Member
			WHERE
				storage_id = $1 AND member_id = $2 AND is_deleted = FALSE
			FOR UPDATE
	`

	var sm model.StorageMember

	err := r.db.QueryRowContext(ctx, query, storageID, userID).Scan(
		&sm.ID,
		&sm.Storage.ID,
		&sm.Member.ID,
		&sm.IsAdmin,
		&sm.IsActive,
		&sm.IsDeleted,
		&sm.CreatedAt,
		&sm.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewError(nil, app.ENotFound)
		}

		return nil, err
	}

	return &sm, nil
}

func (r *repository) FindByStorageID(ctx context.Context, storageID int64) ([]*model.StorageMember, error) {
	query := `
			SELECT
				sm.id, sm.storage_id, p.id, p.photo, p.username, sm.is_admin, sm.is_active, sm.is_deleted, sm.created_at,
				sm.updated_at
			FROM
				Storage_Member sm
			JOIN
				Profile p
			ON
				sm.member_id = p.id
			WHERE
				sm.storage_id = $1 AND sm.is_deleted = FALSE AND p.is_deleted = FALSE
			ORDER BY
				sm.is_admin DESC
	`

	rows, err := r.db.QueryContext(ctx, query, storageID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sMembers := make([]*model.StorageMember, 0)

	for rows.Next() {
		var sm model.StorageMember

		err := rows.Scan(
			&sm.ID,
			&sm.Storage.ID,
			&sm.Member.ID,
			&sm.Member.Photo,
			&sm.Member.Username,
			&sm.IsAdmin,
			&sm.IsActive,
			&sm.IsDeleted,
			&sm.CreatedAt,
			&sm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sMembers = append(sMembers, &sm)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return sMembers, nil
}

func (r *repository) FindByUserID(ctx context.Context, memberID int64) ([]*model.StorageMember, error) {
	query := `
			SELECT
				sm.id, s.id, s.name, s.description, sm.member_id, sm.is_admin, sm.is_active, sm.is_deleted,
				sm.created_at, sm.updated_at
			FROM
				Storage_Member sm
			JOIN
				Storage s
			ON
				sm.storage_id = s.id
			WHERE
				sm.member_id = $1 AND sm.is_deleted = FALSE AND s.is_deleted = FALSE
	`

	rows, err := r.db.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sMembers := make([]*model.StorageMember, 0)

	for rows.Next() {
		var sm model.StorageMember

		err := rows.Scan(
			&sm.ID,
			&sm.Storage.ID,
			&sm.Storage.Name,
			&sm.Storage.Description,
			&sm.Member.ID,
			&sm.IsAdmin,
			&sm.IsActive,
			&sm.IsDeleted,
			&sm.CreatedAt,
			&sm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		sMembers = append(sMembers, &sm)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return sMembers, nil
}

func (r *repository) Update(ctx context.Context, sm *model.StorageMember) error {
	query := `
			UPDATE
				Storage_Member
			SET
				is_admin = $1,
				is_active = $2,
				updated_at = $3
			WHERE
				id = $4
	`

	res, err := r.db.ExecContext(ctx, query, sm.IsAdmin, sm.IsActive, sm.UpdatedAt, sm.ID)
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

func (r *repository) SoftDelete(ctx context.Context, storMembID int64, isDeleted bool) error {
	query := `
			UPDATE
				Storage_Member
			SET
				is_deleted = $1
			WHERE
				id = $2
	`

	res, err := r.db.ExecContext(ctx, query, isDeleted, storMembID)
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
