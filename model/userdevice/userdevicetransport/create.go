package userdevicetransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/userdevice/userdevicebiz"
	"salon_be/model/userdevice/userdevicemodel"
	"salon_be/model/userdevice/userdevicerepo"
	"salon_be/model/userdevice/userdevicestore"

	"github.com/gin-gonic/gin"
)

func CreateUserDevice(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input userdevicemodel.CreateUserDevice

		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("user not found")))
		}

		input.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		store := userdevicestore.NewSQLStore(db)
		repo := userdevicerepo.NewCreateUserDeviceRepo(store)
		biz := userdevicebiz.NewCreateUserDeviceBiz(repo)

		if err := biz.CreateUserDevice(c.Request.Context(), &input); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
