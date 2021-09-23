package utils

import (
	"context"

	"github.com/bagus2x/tjiwi/app"
	"github.com/gin-gonic/gin"
)

type GinCtxKey struct{}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinCtxKey{})
	if ginContext == nil {

		return nil, app.NewError(nil, app.EInternal, "Could not retrieve gin.Context")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, app.NewError(nil, app.EInternal, "gin.Context has wrong type")
	}

	return gc, nil
}
