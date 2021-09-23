package service

import (
	"context"

	"github.com/bagus2x/tjiwi/pkg/history"
)

type service struct {
	historyRepo history.Repository
}

func New(historyRepo history.Repository) history.Service {
	return &service{
		historyRepo: historyRepo,
	}
}

func (s *service) Filter(ctx context.Context, params *history.Params) (*history.GetHistoriesResponse, error) {
	histories, cursor, err := s.historyRepo.Filter(ctx, params)
	if err != nil {
		return nil, err
	}

	var res history.GetHistoriesResponse
	res.Cursor = *cursor
	res.Histories = make([]*history.GetHistoryResponse, 0)

	for _, h := range histories {
		res.Histories = append(res.Histories, &history.GetHistoryResponse{
			ID:        h.ID,
			StorageID: h.Storage.ID,
			BasePaper: history.BasePaper{
				ID:             h.BasePaper.ID,
				Gsm:            h.BasePaper.Gsm,
				Width:          h.BasePaper.Width,
				Io:             h.BasePaper.Io,
				MaterialNumber: h.BasePaper.MaterialNumber,
				Quantity:       h.BasePaper.Quantity,
				Location:       h.BasePaper.Location,
			},
			Member: history.Member{
				ID:       h.Member.ID,
				Photo:    h.Member.Photo.String,
				Username: h.Member.Username,
			},
			Status:    h.Status,
			Affected:  h.Affected,
			CreatedAt: h.CreatedAt,
		})
	}

	return &res, nil
}
