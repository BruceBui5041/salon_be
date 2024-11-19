package bookingtransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	models "salon_be/model"
	"salon_be/model/booking/bookingbiz"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/booking/bookingrepo"
	"salon_be/model/booking/bookingstore"
	"salon_be/model/payment/paymentstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CompleteBookingHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data bookingmodel.CompleteBooking

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrNoPermission(errors.New("requester not found")))
		}

		user, ok := requester.(*models.User)
		if !ok {
			panic(common.ErrNoPermission(errors.New("invalid user type")))
		}

		data.UserID = requester.GetUserId()
		data.IsUserRole = user.IsUser()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			bookingStore := bookingstore.NewSQLStore(tx)
			paymentStore := paymentstore.NewSQLStore(tx)

			repo := bookingrepo.NewCompleteBookingRepo(bookingStore, paymentStore)
			business := bookingbiz.NewCompleteBookingBiz(repo)

			if err := business.CompleteBooking(c.Request.Context(), uid.GetLocalID(), &data); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
