package storagemember

import (
	"context"

	"github.com/bagus2x/tjiwi/pkg/model"
)

type Repository interface {
	Create(ctx context.Context, sm *model.StorageMember) error
	FindByID(ctx context.Context, stormembID int64) (*model.StorageMember, error)
	FindByStorageID(ctx context.Context, storageID int64) ([]*model.StorageMember, error)
	FindByStorageIDAndUserID(ctx context.Context, storageID, userID int64) (*model.StorageMember, error)
	FindByUserID(ctx context.Context, memberID int64) ([]*model.StorageMember, error)
	Update(ctx context.Context, sm *model.StorageMember) error
	SoftDelete(ctx context.Context, stormembID int64, isDeleted bool) error
}
