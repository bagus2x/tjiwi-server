package handler

import (
	"strconv"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/app/middleware"
	"github.com/bagus2x/tjiwi/pkg/storagemember"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StorageMember(r *gin.RouterGroup, service storagemember.Service, mw *middleware.Middleware) {
	r.POST("", mw.AuthJWT(), createStorageMember(service))
	r.GET("/:storMembID", mw.AuthJWT(), getStorageMemberByID(service))
	r.GET("/storage/:storageID", mw.AuthJWT(), getStorageMembersByStorageID(service))
	r.GET("/member/:userID", mw.AuthJWT(), getStorageMembersByUserID(service))
	r.PATCH("/:storMembID", mw.AuthJWT(), updateStorageMember(service))
	r.DELETE("/:storMembID", mw.AuthJWT(), deleteStorageMember(service))
}

func createStorageMember(service storagemember.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req storagemember.CreateStorMembRequest

		err := c.Bind(&req)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid json format"},
				},
			})
			return
		}

		res, err := service.Create(c.Request.Context(), &req)
		if err != nil {
			logrus.Error(err)
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.JSON(201, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func getStorageMemberByID(service storagemember.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storageID := c.Param("storMembID")
		sID, err := strconv.ParseInt(storageID, 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid storage member id"},
				},
			})
			return
		}

		res, err := service.GetByID(c.Request.Context(), sID)
		if err != nil {
			logrus.Error(err)
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

func getStorageMembersByStorageID(service storagemember.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storageID := c.Param("storageID")
		sID, err := strconv.ParseInt(storageID, 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid storage id"},
				},
			})
			return
		}

		res, err := service.GetByStorageID(c.Request.Context(), sID)
		if err != nil {
			logrus.Error(err)
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

func getStorageMembersByUserID(service storagemember.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storageID := c.Param("userID")
		sID, err := strconv.ParseInt(storageID, 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid user id"},
				},
			})
			return
		}

		res, err := service.GetByUserID(c.Request.Context(), sID)
		if err != nil {
			logrus.Error(err)
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

func updateStorageMember(service storagemember.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storMembID := c.Param("storMembID")
		sID, err := strconv.ParseInt(storMembID, 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid storage member id"},
				},
			})
			return
		}

		var req storagemember.UpdateStorMembRequest

		req.ID = sID

		err = c.Bind(&req)
		if err != nil {
			c.JSON(400, app.Failure{
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
			logrus.Error(err)
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

func deleteStorageMember(service storagemember.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storMembID := c.Param("storMembID")
		sID, err := strconv.ParseInt(storMembID, 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid storage member id"},
				},
			})
			return
		}

		err = service.Delete(c.Request.Context(), sID)
		if err != nil {
			logrus.Error(err)
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.ErrorCode(err),
					Messages: app.ErrorMessage(err),
				},
			})
			return
		}

		c.Status(204)
	}
}
