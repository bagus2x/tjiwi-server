package service

import (
	"context"
	"testing"

	"github.com/bagus2x/tjiwi/config"
	"github.com/bagus2x/tjiwi/db"
	"github.com/bagus2x/tjiwi/pkg/user"
	"github.com/bagus2x/tjiwi/pkg/user/repository"
	"github.com/stretchr/testify/assert"
)

var cfg = config.NewTest()

var dbTest = db.NewPostgresDatabase(cfg)

func TestSignUp(t *testing.T) {
	service := New(repository.New(dbTest), cfg)
	res, err := service.SignUp(context.Background(), &user.SignUpRequest{
		Username: "bagus",
		Password: "bagus123",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Log(res)
}

func TestSignIn(t *testing.T) {
	service := New(repository.New(dbTest), cfg)
	res, err := service.SignIn(context.Background(), &user.SignInRequest{
		UsernameOrEmail: "bagus",
		Password:        "bagus123",
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)
	t.Log(res)
}

func TestParseAccessToken(t *testing.T) {
	service := New(repository.New(dbTest), cfg)
	res, err := service.ExtractAccessToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjYzOTg2NDEsIlByb2ZpbGVJRCI6MX0.6DG9L2DOUNQmrBBjNMu2OkZ9ZBabYM-FYZwUQmVccM4")
	assert.NotNil(t, res)
	assert.NoError(t, err)
}
