package service

import (
	"context"
	"time"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/bagus2x/tjiwi/pkg/storage"
	stormemb "github.com/bagus2x/tjiwi/pkg/storagemember"
)

type service struct {
	storageRepo  storage.Repository
	storMembRepo stormemb.Repository
}

func New(storageRepo storage.Repository, storMembRepo stormemb.Repository, cfg *config.Config) storage.Service {
	return &service{
		storageRepo:  storageRepo,
		storMembRepo: storMembRepo,
	}
}

func (s *service) Create(ctx context.Context, req *storage.CreateStorageRequest) (*storage.CreateStorageResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	st := model.Storage{
		Supervisor: model.User{
			ID: req.SupervisorID,
		},
		Name:        req.Name,
		Description: db.NewNullString(req.Description, true),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	err = s.storageRepo.WithTransaction(ctx, func(c context.Context) error {
		err = s.storageRepo.Create(c, &st)
		if err != nil {
			return err
		}

		err := s.storMembRepo.Create(c, &model.StorageMember{
			Storage:   model.Storage{ID: st.ID},
			Member:    model.User{ID: req.SupervisorID},
			IsAdmin:   true,
			IsActive:  true,
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	res := storage.CreateStorageResponse{
		ID:           st.ID,
		SupervisorID: req.SupervisorID,
		Name:         st.Name,
		Description:  st.Description.String,
		CreatedAt:    st.CreatedAt,
		UpdatedAt:    st.UpdatedAt,
	}

	return &res, nil
}

func (s *service) GetByID(ctx context.Context, storageID int64) (*storage.FindStorageResponse, error) {
	st, err := s.storageRepo.FindByID(ctx, storageID)
	if app.ErrorCode(err) == app.ENotFound {
		return nil, app.NewError(err, app.ENotFound, "Storage not found")
	} else if err != nil {
		return nil, err
	}

	res := storage.FindStorageResponse{
		ID: st.ID,
		Supervisor: storage.Supervisor{
			ID:       st.Supervisor.ID,
			Photo:    st.Supervisor.Photo.String,
			Username: st.Supervisor.Username,
			Email:    st.Supervisor.Email,
		},
		Name:        st.Name,
		Description: st.Description.String,
		CreatedAt:   st.CreatedAt,
		UpdatedAt:   st.UpdatedAt,
	}

	return &res, nil
}

func (s *service) GetBySupervisorID(ctx context.Context, supervisorID int64) ([]*storage.FindStorageResponse, error) {
	st, err := s.storageRepo.FindBySupervisorID(ctx, supervisorID)
	if err != nil {
		return nil, err
	}

	res := make([]*storage.FindStorageResponse, 0)

	for _, s := range st {
		res = append(res, &storage.FindStorageResponse{
			ID: s.ID,
			Supervisor: storage.Supervisor{
				ID: s.Supervisor.ID,
			},
			Name:        s.Name,
			Description: s.Description.String,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		})
	}

	return res, nil
}

func (s *service) Update(ctx context.Context, req *storage.UpdateStorageRequest) (*storage.UpdateStorageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	st := model.Storage{
		Name:        req.Name,
		Description: db.NewNullString(req.Description, req.Description != ""),
		UpdatedAt:   time.Now().Unix(),
	}

	err := s.storageRepo.Update(ctx, &st)
	if err != nil {
		return nil, err
	}

	res := storage.UpdateStorageResponse{
		ID:          st.ID,
		Name:        st.Name,
		Description: st.Description.String,
		UpdatedAt:   st.UpdatedAt,
	}

	return &res, nil
}

func (s *service) Delete(ctx context.Context, storageID int64) error {
	err := s.storageRepo.SoftDelete(ctx, storageID, true)
	if app.ErrorCode(err) == app.ENotFound {
		return app.NewError(err, app.ENotFound, "Storage does not exist")
	}

	return err
}
