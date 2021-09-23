package basepaper

import "context"

type Service interface {
	StoreBasePaper(ctx context.Context, req *AddBasePaperRequest) (*AddBasePaperResponse, error)
	GetByID(ctx context.Context, basePaperID int64) (*GetBasePaperResponse, error)
	SearchInBufferArea(ctx context.Context, params *Params) (*GetBasePapersResponse, error)
	SearchInList(ctx context.Context, params *Params) (*GetBasePapersResponse, error)
	MoveToList(ctx context.Context, req *MoveToStorageRequest) (*MoveToStorageResponse, error)
	Deliver(ctx context.Context, req *DeliverBasePaperRequest) (*DeliverBasePaperResponse, error)
	Delete(ctx context.Context, basePaperID int64) error
}
