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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AcceptBookingHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrNoPermission(errors.New("requester not found")))
		}

		// Verify the requester is a provider
		if !requester.IsProvider() {
			panic(common.ErrNoPermission(errors.New("only providers can accept bookings")))
		}

		data := bookingmodel.AcceptBooking{
			UserID: requester.GetUserId(),
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			store := bookingstore.NewSQLStore(tx)
			repo := bookingrepo.NewAcceptBookingRepo(store)
			business := bookingbiz.NewAcceptBookingBiz(repo)

			if err := business.AcceptBooking(c.Request.Context(), uid.GetLocalID(), &data); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
