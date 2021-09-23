package handler

import (
	"strconv"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/app/middleware"
	"github.com/bagus2x/tjiwi/pkg/basepaper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func BasePaper(r *gin.RouterGroup, service basepaper.Service, mw *middleware.Middleware) {
	r.PUT("", mw.AuthJWT(), mw.MustBeStorageMember(false, true), addBasePaper(service))
	r.GET("/:basePaperID", mw.AuthJWT(), mw.MustBeStorageMember(false, true), getBasePaper(service))
	r.GET("/storage/:storageID/search-in-buffer-area", mw.AuthJWT(), mw.MustBeStorageMember(false, true), searchInBufferArea(service))
	r.GET("/storage/:storageID/search-in-list", mw.AuthJWT(), mw.MustBeStorageMember(false, true), searchInList(service))
	r.PUT("/:basePaperID/move-to-list", mw.AuthJWT(), mw.MustBeStorageMember(false, true), moveToList(service))
	r.PUT("/:basePaperID/deliver", mw.AuthJWT(), mw.MustBeStorageMember(false, true), deliver(service))
	r.DELETE("/:basePaperID", mw.AuthJWT(), mw.MustBeStorageMember(true, true), deleteBasePaper(service))
}

func addBasePaper(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req basepaper.AddBasePaperRequest

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

		res, err := service.StoreBasePaper(c.Request.Context(), &req)
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

func getBasePaper(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storageID := c.Param("basePaperID")
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

func searchInBufferArea(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params basepaper.Params
		err := c.Bind(&params)
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

		storageID, err := strconv.ParseInt(c.Param("storageID"), 10, 64)
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

		params.StorageID = &storageID

		res, err := service.SearchInBufferArea(c.Request.Context(), &params)
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

func searchInList(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params basepaper.Params
		err := c.Bind(&params)
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

		storageID, err := strconv.ParseInt(c.Param("storageID"), 10, 64)
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

		params.StorageID = &storageID

		res, err := service.SearchInList(c.Request.Context(), &params)
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

func deliver(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		basePaperID, err := strconv.ParseInt(c.Param("basePaperID"), 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid base paper id"},
				},
			})
			return
		}

		var req basepaper.DeliverBasePaperRequest

		req.ID = basePaperID

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

		res, err := service.Deliver(c.Request.Context(), &req)
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

func moveToList(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		basePaperID, err := strconv.ParseInt(c.Param("basePaperID"), 10, 64)
		if err != nil {
			c.JSON(400, app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
					Messages: []string{"Invalid base paper id"},
				},
			})
			return
		}

		var req basepaper.MoveToStorageRequest

		req.ID = basePaperID

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

		res, err := service.MoveToList(c.Request.Context(), &req)
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

func deleteBasePaper(service basepaper.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		storageID := c.Param("basePaperID")
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

		c.Status(204)
	}
}
