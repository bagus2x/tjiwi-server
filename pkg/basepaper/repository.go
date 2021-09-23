package basepaper

import (
	"context"

	"github.com/bagus2x/tjiwi/pkg/model"
)

type Repository interface {
	Create(ctx context.Context, bp *model.BasePaper) error
	Upsert(ctx context.Context, bp *model.BasePaper) error
	FindByID(ctx context.Context, basePaperID int64) (*model.BasePaper, error)
	Filter(ctx context.Context, params *Params, locationEmpty bool) ([]*model.BasePaper, *Cursor, error)
	Update(ctx context.Context, basePaper *model.BasePaper) error
	SoftDelete(ctx context.Context, basePaperID int64) error
	WithTransaction(ctx context.Context, fn func(context.Context) error) error
}
