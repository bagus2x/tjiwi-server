package utils

import (
	"context"

	"github.com/bagus2x/tjiwi/app"
)

func GetUserIDFromCtx(ctx context.Context) (int64, error) {
	ginCtx, err := GinContextFromContext(ctx)
	if err != nil {
		return 0, err
	}

	userID, ok := ginCtx.Value("userID").(int64)
	if !ok {
		return 0, app.NewError(nil, app.EUnauthorized)
	}

	return userID, nil
}
