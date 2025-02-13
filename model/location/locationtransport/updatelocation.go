package locationtransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	models "salon_be/model"
	"salon_be/model/location/locationbiz"
	"salon_be/model/location/locationmodel"
	"salon_be/model/location/locationrepo"
	"salon_be/model/location/locationstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateLocation(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input locationmodel.UpdateLocationInput
		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("requester not found")))
		}
		input.UserId = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			locationStore := locationstore.NewSQLStore(tx)
			repo := locationrepo.NewUpdateLocationRepo(locationStore)
			biz := locationbiz.NewUpdateLocationBiz(repo)

			location := &models.Location{
				UserId:    input.UserId,
				Latitude:  input.Latitude,
				Longitude: input.Longitude,
				Accuracy:  input.Accuracy,
			}

			if err := biz.UpdateLocation(c.Request.Context(), location); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
