package handler

import (
	"strconv"

	"github.com/bagus2x/tjiwi/app"
	"github.com/bagus2x/tjiwi/app/middleware"
	"github.com/bagus2x/tjiwi/pkg/history"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func History(r *gin.RouterGroup, service history.Service, mw *middleware.Middleware) {
	r.GET("/storage/:storageID/search", searchHistories(service))
}

func searchHistories(service history.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params history.Params
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

		params.StorageID = storageID

		res, err := service.Filter(c.Request.Context(), &params)
		if err != nil {
			logrus.Error(err)
			c.JSON(app.Status(err), app.Failure{
				Success: false,
				Error: app.ErrorDetail{
					Code:     app.EBadRequest,
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
