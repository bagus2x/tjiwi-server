package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/basepaper"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/sirupsen/logrus"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) basepaper.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, bp *model.BasePaper) error {
	query := `
			INSERT INTO
				Base_Paper
				(storage_id, gsm, width, io, material_number, quantity, location, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING
				id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		&bp.Storage.ID,
		&bp.Gsm,
		&bp.Width,
		&bp.Io,
		&bp.MaterialNumber,
		&bp.Quantity,
		&bp.Location,
		&bp.CreatedAt,
		&bp.UpdatedAt,
	).Scan(&bp.ID)

	return err
}

func (r *repository) Upsert(ctx context.Context, bp *model.BasePaper) error {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			INSERT INTO
				Base_Paper
				(storage_id, gsm, width, io, material_number, quantity, location, is_deleted, created_at, updated_at)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			ON CONFLICT
				(storage_id, gsm, width, io, material_number, location)
			DO UPDATE SET
				quantity = Base_Paper.quantity + $6, 
				updated_at = $9, 
				is_deleted = FALSE
			RETURNING
				id, quantity, updated_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		&bp.Storage.ID,
		&bp.Gsm,
		&bp.Width,
		&bp.Io,
		&bp.MaterialNumber,
		&bp.Quantity,
		&bp.Location,
		&bp.IsDeleted,
		&bp.CreatedAt,
		&bp.UpdatedAt,
	).Scan(
		&bp.ID,
		&bp.Quantity,
		&bp.UpdatedAt,
	)

	return err
}

func (r *repository) FindByID(ctx context.Context, basePaperID int64) (*model.BasePaper, error) {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			SELECT
				id, storage_id, gsm, width, io, material_number, quantity, location, created_at, updated_at
			FROM
				Base_Paper
			WHERE
				id = $1 AND is_deleted = FALSE
			FOR UPDATE
	`

	var bp model.BasePaper

	err := tx.QueryRowContext(ctx, query, basePaperID).Scan(
		&bp.ID,
		&bp.Storage.ID,
		&bp.Gsm,
		&bp.Width,
		&bp.Io,
		&bp.MaterialNumber,
		&bp.Quantity,
		&bp.Location,
		&bp.CreatedAt,
		&bp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.NewError(err, app.ENotFound)
		}
		return nil, err
	}

	return &bp, nil
}

func (r *repository) Filter(ctx context.Context, params *basepaper.Params, isLocationEmpty bool) ([]*model.BasePaper, *basepaper.Cursor, error) {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			SELECT
				id, storage_id, gsm, width, io, material_number, quantity, location, created_at, updated_at
			From
				Base_Paper
	`
	dynamicWhere, values := dynamicWhereClause(params, 1, " AND ")
	currentCursor, limit, direction := getCursor(params)
	location := "location != ''"

	if isLocationEmpty {
		location = "location = ''"
	}

	if direction == "next" {
		query = fmt.Sprintf(`
			
			%s
			WHERE
				id > %d AND %s AND %s AND is_deleted = FALSE AND quantity > 0
			ORDER BY
				id ASC
			LIMIT
				%d
			`,
			query, currentCursor, dynamicWhere, location, limit,
		)
	} else {
		query = fmt.Sprintf(`
			WITH prev_mode AS (
				%s
				WHERE
					id < %d AND %s AND %s AND is_deleted = FALSE
				ORDER BY
					id DESC
				LIMIT
					%d
			)
			SELECT
				*
			FROM
				prev_mode 
			ORDER BY 
				id ASC
			`, query, currentCursor, dynamicWhere, location, limit,
		)
	}

	rows, err := tx.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	basePapers := make([]*model.BasePaper, 0)

	for rows.Next() {
		var bp model.BasePaper
		err := rows.Scan(
			&bp.ID,
			&bp.Storage.ID,
			&bp.Gsm,
			&bp.Width,
			&bp.Io,
			&bp.MaterialNumber,
			&bp.Quantity,
			&bp.Location,
			&bp.CreatedAt,
			&bp.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		basePapers = append(basePapers, &bp)
	}

	var cursor basepaper.Cursor

	if len(basePapers) > 0 {
		cursor.Next = basePapers[len(basePapers)-1].ID
		cursor.Previous = basePapers[0].ID
	}

	return basePapers, &cursor, nil
}

func (r *repository) Update(ctx context.Context, basePaper *model.BasePaper) error {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			UPDATE
				Base_Paper
			SET
				gsm = $1, width = $2, io = $3, material_number = $4, quantity = $5,
				location = $6, updated_at = $7
			WHERE
				id = $8 AND is_deleted = FALSE
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		&basePaper.Gsm,
		&basePaper.Width,
		&basePaper.Io,
		&basePaper.MaterialNumber,
		&basePaper.Quantity,
		&basePaper.Location,
		&basePaper.UpdatedAt,
		&basePaper.ID,
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

func (r *repository) SoftDelete(ctx context.Context, basePaperID int64) error {
	tx := db.AllowTransaction(r.db, ctx)

	query := `
			UPDATE
				Base_Paper
			SET
				is_deleted = TRUE,
				quantity = 0
			WHERE
				id = $1
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		basePaperID,
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

func dynamicWhereClause(params *basepaper.Params, index int, sep string) (string, []interface{}) {
	columns := make([]string, 0)
	values := make([]interface{}, 0)

	if params.StorageID != nil {
		columns = append(columns, fmt.Sprintf("%s=$%d", "storage_id", index))
		values = append(values, params.StorageID)

		index++
	}
	if params.Gsm != nil {
		columns = append(columns, fmt.Sprintf("%s=$%d", "gsm", index))
		values = append(values, params.Gsm)

		index++
	}
	if params.Width != nil {
		columns = append(columns, fmt.Sprintf("%s=$%d", "width", index))
		values = append(values, params.Width)
		index++
	}
	if params.Io != nil {
		columns = append(columns, fmt.Sprintf("%s=$%d", "io", index))
		values = append(values, params.Io)
		index++
	}
	if params.MaterialNumber != nil {
		columns = append(columns, fmt.Sprintf("%s=$%d", "material_number", index))
		values = append(values, params.MaterialNumber)
		index++
	}
	if params.Location != nil {
		values = append(values, params.Location)
		columns = append(columns, fmt.Sprintf("%s=$%d", "location", index))
	}

	return strings.Join(columns, " AND "), values
}

func getCursor(params *basepaper.Params) (int64, int64, string) {
	cursor := int64(0)
	limit := int64(10)
	direction := "next"

	if params.NextCursor != nil {
		cursor = *params.NextCursor
	}
	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Direction != nil {
		direction = *params.Direction
	}

	return cursor, limit, direction
}
