package service

import (
	"context"
	"testing"

	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/storage"
	"github.com/bagus2x/tjiwi/pkg/storage/repository"
	"github.com/stretchr/testify/assert"
)

var cfg = config.NewTest()

var dbTest = db.NewPostgresDatabase(cfg)

func TestCreateStorage(t *testing.T) {
	service := New(repository.New(dbTest), nil, cfg)
	res, err := service.Create(context.Background(), &storage.CreateStorageRequest{
		SupervisorID: -1,
		Name:         "Gudang Baru",
		Description:  "Mantap ini",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Log(res)
}

func TestFindStorageByID(t *testing.T) {
	service := New(repository.New(dbTest), nil, cfg)
	res, err := service.GetByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Log(res)
}

func TestFindStorageBySupervisorID(t *testing.T) {
	service := New(repository.New(dbTest), nil, cfg)
	res, err := service.GetBySupervisorID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Log(res)
	t.Log(len(res))
}
