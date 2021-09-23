package service

import (
	"context"
	"strings"
	"time"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/pkg/basepaper"
	"github.com/bagus2x/tjiwi/pkg/history"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/bagus2x/tjiwi/utils"
	"github.com/sirupsen/logrus"
)

type service struct {
	basePaperRepo basepaper.Repository
	historyRepo   history.Repository
}

func New(basePaperRepo basepaper.Repository, historyRepo history.Repository) basepaper.Service {
	return &service{
		basePaperRepo: basePaperRepo,
		historyRepo:   historyRepo,
	}
}

func (s *service) StoreBasePaper(ctx context.Context, req *basepaper.AddBasePaperRequest) (*basepaper.AddBasePaperResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var bp model.BasePaper

	err := s.basePaperRepo.WithTransaction(ctx, func(c context.Context) error {
		bp = model.BasePaper{
			Storage:        model.Storage{ID: req.StorageID},
			Gsm:            req.Gsm,
			Width:          req.Width,
			Io:             req.Io,
			MaterialNumber: req.MaterialNumber,
			Quantity:       req.Quantity,
			CreatedAt:      time.Now().Unix(),
			UpdatedAt:      time.Now().Unix(),
		}

		err := s.basePaperRepo.Upsert(c, &bp)
		if err != nil {
			logrus.Error("error upsert")
			return app.NewError(err, app.EInternal, "Failed to save base paper")
		}

		memberID, err := utils.GetUserIDFromCtx(c)
		if err != nil {
			return err
		}

		err = s.historyRepo.Create(c, &model.History{
			BasePaper: model.BasePaper{ID: bp.ID},
			Storage:   bp.Storage,
			Member:    model.User{ID: memberID},
			Status:    "stored",
			Affected:  bp.Quantity,
			CreatedAt: bp.UpdatedAt,
		})
		if err != nil {
			logrus.Error("error create history")
		}

		return err

	})
	if err != nil {
		return nil, err
	}

	res := basepaper.AddBasePaperResponse{
		ID:             bp.ID,
		StorageID:      bp.Storage.ID,
		Gsm:            bp.Gsm,
		Width:          bp.Width,
		Io:             bp.Io,
		MaterialNumber: bp.MaterialNumber,
		Quantity:       bp.Quantity,
		CreatedAt:      bp.CreatedAt,
		UpdatedAt:      bp.UpdatedAt,
	}

	return &res, nil
}

func (s *service) GetByID(ctx context.Context, basePaperID int64) (*basepaper.GetBasePaperResponse, error) {
	bp, err := s.basePaperRepo.FindByID(ctx, basePaperID)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(nil, app.ENotFound)
	} else if err != nil {
		return nil, err
	}

	res := basepaper.GetBasePaperResponse{
		ID:             bp.ID,
		StorageID:      bp.Storage.ID,
		Gsm:            bp.Gsm,
		Width:          bp.Width,
		Io:             bp.Io,
		MaterialNumber: bp.MaterialNumber,
		Location:       bp.Location,
		Quantity:       bp.Quantity,
		CreatedAt:      bp.CreatedAt,
		UpdatedAt:      bp.UpdatedAt,
	}

	return &res, nil
}

func (s *service) SearchInBufferArea(ctx context.Context, params *basepaper.Params) (*basepaper.GetBasePapersResponse, error) {
	return s.search(ctx, params, true)
}

func (s *service) SearchInList(ctx context.Context, params *basepaper.Params) (*basepaper.GetBasePapersResponse, error) {
	return s.search(ctx, params, false)
}

func (s *service) search(ctx context.Context, params *basepaper.Params, isLocationEmpty bool) (*basepaper.GetBasePapersResponse, error) {
	basepapers, cursor, err := s.basePaperRepo.Filter(ctx, params, isLocationEmpty)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(nil, app.ENotFound)
	} else if err != nil {
		return nil, err
	}

	basepapersRes := make([]*basepaper.GetBasePaperResponse, 0)
	for _, bp := range basepapers {
		basepapersRes = append(basepapersRes, &basepaper.GetBasePaperResponse{
			ID:             bp.ID,
			StorageID:      bp.Storage.ID,
			Gsm:            bp.Gsm,
			Width:          bp.Width,
			Io:             bp.Io,
			MaterialNumber: bp.MaterialNumber,
			Location:       bp.Location,
			Quantity:       bp.Quantity,
			CreatedAt:      bp.CreatedAt,
			UpdatedAt:      bp.UpdatedAt,
		})
	}

	res := basepaper.GetBasePapersResponse{
		Cursor:     *cursor,
		BasePapers: basepapersRes,
	}

	return &res, nil
}

func (s *service) MoveToList(ctx context.Context, req *basepaper.MoveToStorageRequest) (*basepaper.MoveToStorageResponse, error) {
	var res basepaper.MoveToStorageResponse

	err := s.basePaperRepo.WithTransaction(ctx, func(c context.Context) error {
		bp, err := s.basePaperRepo.FindByID(c, req.ID)
		if app.ErrorCode(err) == app.ENotFound {
			return app.NewError(nil, app.ENotFound, "Base paper not found")
		} else if err != nil {
			return err
		}
		if bp.Location != "" || bp.Quantity == 0 {
			return app.NewError(nil, app.ENotFound, "Base paper not found")
		}
		if bp.Quantity-req.Quantity < 0 {
			return app.NewError(nil, app.EBadRequest, "Quantity exceeds the limit")
		}

		bp.Quantity -= req.Quantity
		bp.UpdatedAt = time.Now().Unix()

		err = s.basePaperRepo.Update(c, bp)
		if app.ErrorCode(err) == app.ENotFound {
			return err
		}

		bp.Quantity = req.Quantity
		bp.Location = strings.ToUpper(req.Location)
		bp.CreatedAt = bp.UpdatedAt

		err = s.basePaperRepo.Upsert(c, bp)
		if err != nil {
			return err
		}

		memberID, err := utils.GetUserIDFromCtx(c)
		if err != nil {
			return err
		}

		err = s.historyRepo.Create(c, &model.History{
			BasePaper: model.BasePaper{ID: bp.ID},
			Storage:   bp.Storage,
			Member:    model.User{ID: memberID},
			Status:    "moved",
			Affected:  req.Quantity,
			CreatedAt: bp.UpdatedAt,
		})
		if err != nil {
			return err
		}

		res.ID = bp.ID
		res.StorageID = bp.Storage.ID
		res.Gsm = bp.Gsm
		res.Width = bp.Width
		res.Io = bp.Io
		res.MaterialNumber = bp.MaterialNumber
		res.Location = bp.Location
		res.Quantity = bp.Quantity
		res.UpdatedAt = bp.UpdatedAt
		res.CreatedAt = bp.CreatedAt

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (s *service) Deliver(ctx context.Context, req *basepaper.DeliverBasePaperRequest) (*basepaper.DeliverBasePaperResponse, error) {
	var res basepaper.DeliverBasePaperResponse
	err := s.basePaperRepo.WithTransaction(ctx, func(c context.Context) error {
		bp, err := s.basePaperRepo.FindByID(c, req.ID)
		if app.ErrorCode(err) == app.ENotFound {
			return app.NewError(nil, app.ENotFound, "Base paper not found")
		} else if err != nil {
			return err
		}
		if bp.Location == "" || bp.Quantity == 0 {
			return app.NewError(nil, app.ENotFound, "Base paper not found")
		}
		if bp.Quantity-req.Quantity < 0 {
			return app.NewError(nil, app.EBadRequest, "Quantity exceeds the limit")
		}

		bp.Quantity -= req.Quantity

		err = s.basePaperRepo.Update(c, bp)
		if err != nil {
			return err
		}

		memberID, err := utils.GetUserIDFromCtx(c)
		if err != nil {
			return err
		}

		err = s.historyRepo.Create(c, &model.History{
			BasePaper: model.BasePaper{ID: bp.ID},
			Storage:   bp.Storage,
			Member:    model.User{ID: memberID},
			Status:    "delivered",
			Affected:  req.Quantity,
			CreatedAt: bp.CreatedAt,
		})
		if err != nil {
			return err
		}

		res = basepaper.DeliverBasePaperResponse(*req)

		return nil
	})

	return &res, err
}

func (s *service) Delete(ctx context.Context, basePaperID int64) error {
	err := s.basePaperRepo.WithTransaction(ctx, func(c context.Context) error {
		bp, err := s.basePaperRepo.FindByID(c, basePaperID)
		if app.ErrorCode(err) == app.ENotFound {
			return app.NewError(err, app.ENotFound, "Base paper not found")
		} else if err != nil {
			return err
		}

		quantity := bp.Quantity

		err = s.basePaperRepo.SoftDelete(c, basePaperID)
		if err != nil {
			return err
		}

		memberID, err := utils.GetUserIDFromCtx(c)
		if err != nil {
			return err
		}

		err = s.historyRepo.Create(c, &model.History{
			BasePaper: model.BasePaper{ID: bp.ID},
			Storage:   bp.Storage,
			Member:    model.User{ID: memberID},
			Status:    "deleted",
			Affected:  quantity,
			CreatedAt: time.Now().Unix(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
