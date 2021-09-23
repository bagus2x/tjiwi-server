package model

import "database/sql"

type Storage struct {
	ID          int64
	Supervisor  User
	Name        string
	Description sql.NullString
	IsDeleted   bool
	CreatedAt   int64
	UpdatedAt   int64
}
