package repository

import (
	"context"
	"testing"

	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/stretchr/testify/assert"
)

var cfg = config.NewTest()
var dbTest = db.NewPostgresDatabase(cfg)

func TestCreateStorageMember(t *testing.T) {
	repo := New(dbTest)
	sm := model.StorageMember{
		Storage: model.Storage{ID: 1},
		Member:  model.User{ID: 1},
	}
	err := repo.Create(context.Background(), &sm)
	assert.NoError(t, err)
}

func TestFindByID(t *testing.T) {
	repo := New(dbTest)
	sm, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, sm)
}

func TestFindByStorageID(t *testing.T) {
	repo := New(dbTest)
	sm, err := repo.FindByStorageID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, sm)
}
