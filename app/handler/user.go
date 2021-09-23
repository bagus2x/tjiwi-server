package handler

import (
	"log"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/app/middleware"
	"github.com/bagus2x/tjiwi/pkg/user"
	"github.com/gin-gonic/gin"
)

func User(r *gin.RouterGroup, service user.Service, mw *middleware.Middleware) {
	r.POST("/signup", signUp(service))
	r.POST("/signin", signIn(service))
	r.POST("/signout", mw.AuthJWT(), signOut(service))
	r.POST("/refresh", refreshToken(service))
	r.GET("/search", searchUsernames(service))
	r.GET("", mw.AuthJWT(), getUser(service))
	r.PUT("", mw.AuthJWT(), updateUser(service))
	r.DELETE("", mw.AuthJWT(), deleteUser(service))
}

func signIn(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.SignInRequest

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
			return
		}

		res, err := service.SignIn(c.Request.Context(), &req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func signUp(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.SignUpRequest

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
			return
		}

		res, err := service.SignUp(c.Request.Context(), &req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func signOut(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		pID, _ := userID.(int64)

		err := service.SignOut(c.Request.Context(), pID)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    userID,
		})
	}
}

func getUser(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		pID, _ := userID.(int64)

		res, err := service.GetUserByID(c.Request.Context(), pID)
		if err != nil {
			log.Println(err)
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func searchUsernames(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")

		res, err := service.SearchByUsername(c.Request.Context(), username)
		if err != nil {
			log.Println(err)
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func updateUser(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		pID, _ := userID.(int64)

		var req user.UpdateUserRequest

		req.ID = pID

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
			return
		}

		res, err := service.Update(c.Request.Context(), &req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func deleteUser(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		pID, _ := userID.(int64)

		err := service.Delete(c.Request.Context(), pID)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    userID,
		})
	}
}

func refreshToken(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req user.RefreshTokenRequest

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
			return
		}

		res, err := service.RefreshToken(c.Request.Context(), &req)
		log.Println(err)
		if err != nil {
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
		}

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}
