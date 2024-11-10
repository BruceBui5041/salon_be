package bookingtransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/booking/bookingbiz"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/booking/bookingrepo"
	"salon_be/model/booking/bookingstore"
	"salon_be/model/serviceversion/serviceversionstore"
	"salon_be/model/user/userstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateBookingHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data bookingmodel.CreateBooking

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrNoPermission(errors.New("requester not found")))
		}
		data.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			bookingStore := bookingstore.NewSQLStore(tx)
			serviceVersionStore := serviceversionstore.NewSQLStore(tx)
			serviceManStore := userstore.NewSQLStore(tx)

			repo := bookingrepo.NewCreateBookingRepo(
				bookingStore,
				serviceVersionStore,
				serviceManStore,
			)

			business := bookingbiz.NewCreateBookingBiz(repo)

			if err := business.CreateBooking(c.Request.Context(), &data); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
