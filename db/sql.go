package db

import (
	"context"
	"database/sql"
)

type Key int

type TransactionKey struct{}

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func AllowTransaction(db *sql.DB, ctx context.Context) Transaction {
	if tx, ok := ctx.Value(TransactionKey{}).(*sql.Tx); ok {
		return tx
	}

	return db
}

func NewNullInt(value int64, valid bool) sql.NullInt64 {
	return sql.NullInt64{
		Int64: value,
		Valid: valid,
	}
}

func NewNullString(value string, valid bool) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  valid,
	}
}
