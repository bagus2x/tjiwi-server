package middleware

import (
	"context"

	"github.com/bagus2x/tjiwi/utils"
	"github.com/gin-gonic/gin"
)

func (mw *Middleware) GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), utils.GinCtxKey{}, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
