package repository

import (
	"context"
	"database/sql"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/bagus2x/tjiwi/pkg/user"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) user.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, p *model.User) error {
	query := `
			INSERT INTO
				Profile
				(photo, username, email, password, is_deleted, token, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING
				id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		p.Photo,
		p.Username,
		p.Email,
		p.Password,
		p.IsDeleted,
		p.Token,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&p.ID)

	return err
}

func (r *repository) FindByID(ctx context.Context, userID int64) (*model.User, error) {
	query := `
			SELECT
				id, photo, username, email, password, token, created_at, updated_at
			FROM
				Profile
			WHERE
				id = $1 AND is_deleted = FALSE
	`

	var p model.User

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&p.ID,
		&p.Photo,
		&p.Username,
		&p.Email,
		&p.Password,
		&p.Token,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewError(err, app.ENotFound)
		}
		return nil, err
	}

	return &p, nil
}

func (r *repository) FindByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error) {
	query := `
			SELECT
				id, photo, username, email, password, created_at, updated_at
			FROM
				Profile
			WHERE
				(username = $1 OR email = $2) AND is_deleted = FALSE
	`

	var p model.User

	err := r.db.QueryRowContext(
		ctx,
		query,
		username,
		email,
	).Scan(
		&p.ID,
		&p.Photo,
		&p.Username,
		&p.Email,
		&p.Password,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewError(err, app.ENotFound)
		}
		return nil, err
	}

	return &p, nil
}

func (r *repository) MatchByUsername(ctx context.Context, username string) ([]*model.User, error) {
	query := `
			SELECT
				id, photo, username, email, password, created_at, updated_at
			FROM
				Profile
			WHERE
				username ILIKE $1 AND is_deleted = FALSE
			LIMIT
				5
	`

	rows, err := r.db.QueryContext(ctx, query, "%"+username+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]*model.User, 0)

	for rows.Next() {
		var u model.User

		rows.Scan(
			&u.ID,
			&u.Photo,
			&u.Username,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.UpdatedAt,
		)

		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, app.NewError(nil, app.ENotFound)
	}

	return users, err
}

func (r *repository) Update(ctx context.Context, p *model.User) error {
	query := `
			UPDATE
				Profile
			SET
				photo = $1,
				username = $2,
				email = $3,
				updated_at = $4
			WHERE
				id = $5
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		p.Photo,
		p.Username,
		p.Email,
		p.UpdatedAt,
		p.ID,
	)
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

func (r *repository) SoftDelete(ctx context.Context, userID int64, isDeleted bool) error {
	query := `
			UPDATE
				Profile
			SET
				is_deleted = $1
			WHERE
				id = $2
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		isDeleted,
		userID,
	)
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

func (r *repository) UpdateToken(ctx context.Context, userID int64, token string) error {
	query := `
			UPDATE
				Profile
			SET
				token = $1
			WHERE
				id = $2
	`
	res, err := r.db.ExecContext(
		ctx,
		query,
		token,
		userID,
	)
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
