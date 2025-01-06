package bookingtransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/model/booking/bookingbiz"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/booking/bookingrepo"
	"salon_be/model/booking/bookingstore"
	"salon_be/model/payment/paymentstore"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CancelBookingHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingUid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data bookingmodel.CancelBooking

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
			paymentStore := paymentstore.NewSQLStore(tx)
			repo := bookingrepo.NewCancelBookingRepo(bookingStore, paymentStore)
			business := bookingbiz.NewCancelBookingBiz(repo)

			if err := business.CancelBooking(c.Request.Context(), bookingUid.GetLocalID(), &data); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		if err := watermill.PublishBookingEvent(
			c.Request.Context(),
			appCtx.GetLocalPubSub().GetUnblockPubSub(),
			&messagemodel.BookingEventMsg{
				BookingID: bookingUid.GetLocalID(),
				Event:     messagemodel.BookingCancelledEvent,
			},
		); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "error publishing booking event", zap.Error(err))
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
