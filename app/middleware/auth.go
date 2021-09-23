package middleware

import (
	"strconv"
	"strings"

	"github.com/bagus2x/tjiwi/app"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		bearer := strings.Split(authHeader, " ")
		if len(bearer) != 2 {
			c.JSON(401, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EUnauthorized,
					Messages: []string{"Invalid authorization header format"},
				},
			})
			c.Abort()
			return
		}

		claims, err := m.userService.ExtractAccessToken(bearer[1])
		if err != nil {
			c.JSON(401, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
	}
}

func (m *Middleware) MustBeSupervisor() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (m *Middleware) MustBeStorageMember(isAdmin, isActive bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		storMembID, err := strconv.ParseInt(c.GetHeader("X-Storage-Member"), 10, 64)
		if err != nil {
			c.JSON(403, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EForbidden,
					Messages: []string{"X-Storage-Member header is required"},
				},
			})
			c.Abort()
			return
		}

		res, err := m.storMembService.GetByID(c.Request.Context(), storMembID)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			c.Abort()
			return
		}

		userIDInterface, _ := c.Get("userID")
		userID, _ := userIDInterface.(int64)

		if res.Member.ID != userID {
			c.JSON(403, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EForbidden,
					Messages: []string{"Access Denied"},
				},
			})
			c.Abort()
			return
		}

		if isActive && !res.IsActive {
			c.JSON(403, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EForbidden,
					Messages: []string{"User status is inactive"},
				},
			})
			c.Abort()
			return
		}

		if isAdmin && !res.IsAdmin {
			c.JSON(403, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EForbidden,
					Messages: []string{"User status is not an admin"},
				},
			})
			c.Abort()
			return
		}

		c.Set("storageMember", res)
	}
}
