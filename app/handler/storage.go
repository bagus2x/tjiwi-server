package handler

import (
	"strconv"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/app/middleware"
	"github.com/bagus2x/tjiwi/pkg/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Storage(r *gin.RouterGroup, service storage.Service, mw *middleware.Middleware) {
	r.POST("", mw.AuthJWT(), createStorage(service))
	r.GET("", mw.AuthJWT(), getStorages(service))
	r.GET("/:storageID", mw.AuthJWT(), getStorage(service))
	r.DELETE("/:storageID", mw.AuthJWT(), deleteStorage(service))
}

func createStorage(service storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		supervisorID, _ := c.Get("userID")
		sID, _ := supervisorID.(int64)

		var req storage.CreateStorageRequest

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

		req.SupervisorID = sID

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

		c.JSON(200, app.Success{
			Success: true,
			Data:    res,
		})
	}
}

func getStorages(service storage.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		supervisorID, _ := c.Get("userID")
		sID, _ := supervisorID.(int64)

		res, err := service.GetBySupervisorID(c.Request.Context(), sID)
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

func getStorage(service storage.Service) gin.HandlerFunc {
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

func deleteStorage(service storage.Service) gin.HandlerFunc {
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

		c.JSON(200, app.Success{
			Success: true,
			Data:    sID,
		})
	}
}
