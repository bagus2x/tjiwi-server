package user

import (
	"context"

	"github.com/bagus2x/tjiwi/pkg/model"
)

type Repository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, userID int64) (*model.User, error)
	FindByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error)
	MatchByUsername(ctx context.Context, username string) ([]*model.User, error)
	Update(ctx context.Context, user *model.User) error
	SoftDelete(ctx context.Context, userID int64, isDeleted bool) error
	UpdateToken(ctx context.Context, userID int64, token string) error
}
