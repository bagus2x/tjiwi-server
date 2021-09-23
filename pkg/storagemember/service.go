package storagemember

import "context"

type Service interface {
	Create(ctx context.Context, req *CreateStorMembRequest) (*CreateStorMembResponse, error)
	GetByID(ctx context.Context, stormembID int64) (*GetStorMembResponse, error)
	GetByStorageIDAndUserID(ctx context.Context, storageID, userID int64) (*GetStorMembResponse, error)
	GetByStorageID(ctx context.Context, stormembID int64) ([]*GetStorMembResponse, error)
	GetByUserID(ctx context.Context, storMembID int64) ([]*GetStorMembResponse, error)
	Update(ctx context.Context, req *UpdateStorMembRequest) (*UpdateStorMembResponse, error)
	Delete(ctx context.Context, storMembID int64) error
}
