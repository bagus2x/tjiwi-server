package repository

import (
	"context"
	"testing"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/model"
	"github.com/stretchr/testify/assert"
)

var cfg = config.NewTest()

var dbTest = db.NewPostgresDatabase(cfg)

func TestCreate(t *testing.T) {
	repo := New(dbTest)
	p := model.User{
		Username: "anonim",
		Password: "anonim123",
	}
	err := repo.Create(context.Background(), &p)
	assert.NoError(t, err)
	assert.NotZero(t, p.ID)
}

func TestFindByID(t *testing.T) {
	repo := New(dbTest)

	p, err := repo.FindByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotZero(t, p.ID)
	t.Log(p)
}

func TestFindByUsername(t *testing.T) {
	repo := New(dbTest)

	p, err := repo.FindByUsernameOrEmail(context.Background(), "anonim", "anonim@gmail.com")
	assert.NoError(t, err)
	t.Log(p)
	t.Log(err)
	t.Log(app.ErrorCode(err))
	t.Log(app.ErrorMessage(err))
}

func TestUpdate(t *testing.T) {
	repo := New(dbTest)

	err := repo.Update(context.Background(), &model.User{
		ID:       1,
		Username: "budi",
		Password: "rahmat",
		Token:    "wkowkdwd",
	})
	assert.NoError(t, err)
}

func TestUpdateToken(t *testing.T) {
	repo := New(dbTest)

	err := repo.UpdateToken(context.Background(), 1, " ini token")
	assert.NoError(t, err)
}

func TestSoftToken(t *testing.T) {
	repo := New(dbTest)

	err := repo.SoftDelete(context.Background(), 1, true)
	assert.NoError(t, err)
}
