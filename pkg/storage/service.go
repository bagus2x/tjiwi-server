package storage

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, req *CreateStorageRequest) (*CreateStorageResponse, error)
	GetByID(ctx context.Context, storageID int64) (*FindStorageResponse, error)
	GetBySupervisorID(ctx context.Context, storageID int64) ([]*FindStorageResponse, error)
	Update(ctx context.Context, req *UpdateStorageRequest) (*UpdateStorageResponse, error)
	Delete(ctx context.Context, storageID int64) error
}
