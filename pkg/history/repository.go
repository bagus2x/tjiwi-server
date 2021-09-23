package history

import (
	"context"

	"github.com/bagus2x/tjiwi/pkg/model"
)

type Repository interface {
	Create(ctx context.Context, history *model.History) error
	Filter(ctx context.Context, params *Params) ([]*model.History, *Cursor, error)
}
