package storage

import (
	"context"

	"github.com/bagus2x/tjiwi/pkg/model"
)

type Repository interface {
	Create(ctx context.Context, st *model.Storage) error
	FindByID(ctx context.Context, storageID int64) (*model.Storage, error)
	FindBySupervisorID(ctx context.Context, supervisorID int64) ([]*model.Storage, error)
	Update(ctx context.Context, st *model.Storage) error
	SoftDelete(ctx context.Context, storageID int64, isDeleted bool) error
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
