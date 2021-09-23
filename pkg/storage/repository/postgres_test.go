package repository

import (
	"context"
	"testing"
	"time"

	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/stretchr/testify/assert"
)

var cfg = config.NewTest()

var dbTest = db.NewPostgresDatabase(cfg)

func TestCreateStorage(t *testing.T) {
	repo := New(dbTest)
	s := model.Storage{
		Supervisor:  model.User{ID: 1},
		Name:        "Gudang B",
		Description: db.NewNullString("Ini gudang", true),
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}
	err := repo.Create(context.Background(), &s)
	assert.NoError(t, err)
	assert.NotZero(t, s.ID)
}

func TestFindStorageByID(t *testing.T) {
	repo := New(dbTest)
	s, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	t.Log(s)
}

func TestFindStorageBySupervisorID(t *testing.T) {
	repo := New(dbTest)
	s, err := repo.FindBySupervisorID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	t.Log(len(s))
	t.Log(s)
}

func TestUpdateStorage(t *testing.T) {
	repo := New(dbTest)
	err := repo.Update(context.Background(), &model.Storage{
		ID:          1,
		Name:        "Gudang Z",
		Description: db.NewNullString("Ini gudang z", true),
		UpdatedAt:   time.Now().Unix(),
	})
	assert.NoError(t, err)
}

func TestSoftDelete(t *testing.T) {
	repo := New(dbTest)
	err := repo.SoftDelete(context.Background(), 1, false)
	assert.NoError(t, err)
}
