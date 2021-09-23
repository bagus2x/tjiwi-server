package history

import "context"

type Service interface {
	Filter(ctx context.Context, params *Params) (*GetHistoriesResponse, error)
}
