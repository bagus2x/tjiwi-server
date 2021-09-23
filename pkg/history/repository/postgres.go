package repository

import (
	"context"
	"database/sql"

	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/history"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/sirupsen/logrus"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) history.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, history *model.History) error {
	tx := db.AllowTransaction(r.db, ctx)

	logrus.Info("StorID ", history.Storage.ID)
	logrus.Info("MemberID ", history.Member.ID)

	query := `
			INSERT INTO
				History
				(base_paper_id, storage_id, member_id, status, affected, created_at)
			VALUES
				($1, $2, $3, $4, $5, $6)
			RETURNING
				id
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		history.BasePaper.ID,
		history.Storage.ID,
		history.Member.ID,
		history.Status,
		history.Affected,
		history.CreatedAt,
	).Scan(&history.ID)

	return err
}

func (r *repository) Filter(ctx context.Context, params *history.Params) ([]*model.History, *history.Cursor, error) {
	query, values := descendingFilter(params)

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, nil, err
	}

	defer rows.Close()

	histories := make([]*model.History, 0)

	for rows.Next() {
		var history model.History

		err := rows.Scan(
			&history.ID,
			&history.BasePaper.ID,
			&history.BasePaper.Gsm,
			&history.BasePaper.Width,
			&history.BasePaper.Io,
			&history.BasePaper.MaterialNumber,
			&history.BasePaper.Quantity,
			&history.BasePaper.Location,
			&history.Storage.ID,
			&history.Member.ID,
			&history.Member.Photo,
			&history.Member.Username,
			&history.Status,
			&history.Affected,
			&history.CreatedAt,
		)
		if err != nil {
			return nil, nil, err
		}

		histories = append(histories, &history)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	var cursor history.Cursor
	if len(histories) > 0 {
		cursor.Next = histories[len(histories)-1].ID
		cursor.Previous = histories[0].ID
	}

	return histories, &cursor, nil
}
