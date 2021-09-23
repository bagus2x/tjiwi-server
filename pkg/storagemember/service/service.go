package service

import (
	"context"
	"time"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/pkg/model"
	stormemb "github.com/bagus2x/tjiwi/pkg/storagemember"
)

type service struct {
	stormembRepo stormemb.Repository
}

func New(stormembRepo stormemb.Repository, cfg *config.Config) stormemb.Service {
	return &service{
		stormembRepo: stormembRepo,
	}
}

func (s *service) Create(ctx context.Context, req *stormemb.CreateStorMembRequest) (*stormemb.CreateStorMembResponse, error) {
	sm := model.StorageMember{
		Storage: model.Storage{
			ID: req.StorageID,
		},
		Member: model.User{
			ID: req.MemberID,
		},
		IsAdmin:   req.IsAdmin,
		IsActive:  req.IsActive,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := s.stormembRepo.Create(ctx, &sm); err != nil {
		return nil, err
	}

	res := stormemb.CreateStorMembResponse{
		ID:        sm.ID,
		StorageID: sm.Storage.ID,
		MemberID:  sm.Member.ID,
		IsAdmin:   sm.IsAdmin,
		IsActive:  sm.IsActive,
		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}

	return &res, nil
}

func (s *service) GetByID(ctx context.Context, storMembID int64) (*stormemb.GetStorMembResponse, error) {
	sm, err := s.stormembRepo.FindByID(ctx, storMembID)
	if err != nil {
		return nil, err
	}

	res := &stormemb.GetStorMembResponse{
		ID: sm.ID,
		Storage: stormemb.Storage{
			ID: sm.Storage.ID,
		},
		Member: stormemb.Member{
			ID:       sm.Member.ID,
			Username: sm.Member.Username,
			Photo:    sm.Member.Photo.String,
		},
		IsAdmin:   sm.IsAdmin,
		IsActive:  sm.IsActive,
		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetByStorageIDAndUserID(ctx context.Context, storageID, userID int64) (*stormemb.GetStorMembResponse, error) {
	sm, err := s.stormembRepo.FindByStorageIDAndUserID(ctx, storageID, userID)
	if err != nil {
		return nil, err
	}

	res := &stormemb.GetStorMembResponse{
		ID: sm.ID,
		Storage: stormemb.Storage{
			ID: sm.Storage.ID,
		},
		Member: stormemb.Member{
			ID:       sm.Member.ID,
			Username: sm.Member.Username,
			Photo:    sm.Member.Photo.String,
		},
		IsAdmin:   sm.IsAdmin,
		IsActive:  sm.IsActive,
		CreatedAt: sm.CreatedAt,
		UpdatedAt: sm.UpdatedAt,
	}

	return res, nil
}

func (s *service) GetByStorageID(ctx context.Context, storMembID int64) ([]*stormemb.GetStorMembResponse, error) {
	sm, err := s.stormembRepo.FindByStorageID(ctx, storMembID)
	if err != nil {
		return nil, err
	}

	res := make([]*stormemb.GetStorMembResponse, 0)

	for _, storMemb := range sm {
		res = append(res, &stormemb.GetStorMembResponse{
			ID: storMemb.ID,
			Storage: stormemb.Storage{
				ID: storMemb.Storage.ID,
			},
			Member: stormemb.Member{
				ID:       storMemb.Member.ID,
				Username: storMemb.Member.Username,
				Photo:    storMemb.Member.Photo.String,
			},
			IsAdmin:   storMemb.IsAdmin,
			IsActive:  storMemb.IsActive,
			CreatedAt: storMemb.CreatedAt,
			UpdatedAt: storMemb.UpdatedAt,
		})
	}

	return res, nil
}

func (s *service) GetByUserID(ctx context.Context, storMembID int64) ([]*stormemb.GetStorMembResponse, error) {
	sm, err := s.stormembRepo.FindByUserID(ctx, storMembID)
	if err != nil {
		return nil, err
	}

	res := make([]*stormemb.GetStorMembResponse, 0)

	for _, storMemb := range sm {
		res = append(res, &stormemb.GetStorMembResponse{
			ID: storMemb.ID,
			Storage: stormemb.Storage{
				ID:          storMemb.Storage.ID,
				Name:        storMemb.Storage.Name,
				Description: storMemb.Storage.Description.String,
			},
			Member: stormemb.Member{
				ID: storMemb.Member.ID,
			},
			IsAdmin:   storMemb.IsAdmin,
			IsActive:  storMemb.IsActive,
			CreatedAt: storMemb.CreatedAt,
			UpdatedAt: storMemb.UpdatedAt,
		})
	}

	return res, nil
}

func (s *service) Update(ctx context.Context, req *stormemb.UpdateStorMembRequest) (*stormemb.UpdateStorMembResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	sm := model.StorageMember{
		ID:        req.ID,
		IsAdmin:   req.IsAdmin,
		IsActive:  req.IsActive,
		UpdatedAt: time.Now().Unix(),
	}

	err := s.stormembRepo.Update(ctx, &sm)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(nil, app.ENotFound, "Storage member not found")
	} else if err != nil {
		return nil, err
	}

	res := stormemb.UpdateStorMembResponse{
		ID:        sm.ID,
		IsAdmin:   sm.IsAdmin,
		IsActive:  sm.IsActive,
		UpdatedAt: sm.UpdatedAt,
	}

	return &res, nil
}

func (s *service) Delete(ctx context.Context, storMembID int64) error {
	err := s.stormembRepo.SoftDelete(ctx, storMembID, true)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(nil, app.ENotFound, "Storage member not found")
	}

	return err
}
