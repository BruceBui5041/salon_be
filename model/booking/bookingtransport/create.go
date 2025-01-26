package bookingtransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/booking/bookingbiz"
	"salon_be/model/booking/bookingmodel"
	"salon_be/model/booking/bookingrepo"
	"salon_be/model/booking/bookingstore"
	"salon_be/model/coupon/couponstore"
	"salon_be/model/payment/paymentstore"
	"salon_be/model/service/servicestore"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

		user, ok := requester.(*models.User)
		if !ok {
			panic(common.ErrNoPermission(errors.New("invalid user type")))
		}

		data.UserID = requester.GetUserId()
		data.IsUserRole = user.IsUser()

		db := appCtx.GetMainDBConnection()
		var newBookingId uint32
		if err := db.Transaction(func(tx *gorm.DB) error {
			bookingStore := bookingstore.NewSQLStore(tx)
			serviceStore := servicestore.NewSQLStore(tx)
			paymentStore := paymentstore.NewSQLStore(tx)
			couponStore := couponstore.NewSQLStore(tx)

			repo := bookingrepo.NewCreateBookingRepo(
				bookingStore,
				serviceStore,
				paymentStore,
				couponStore,
			)

			business := bookingbiz.NewCreateBookingBiz(repo)

			id, err := business.CreateBooking(c.Request.Context(), &data)
			if err != nil {
				return err
			}

			newBookingId = id

			return nil
		}); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "error creating booking", zap.Error(err))
			panic(err)
		}

		if err := watermill.PublishBookingEvent(
			c.Request.Context(),
			appCtx.GetLocalPubSub().GetUnblockPubSub(),
			&messagemodel.BookingEventMsg{
				BookingID: newBookingId,
				Event:     messagemodel.BookingCreatedEvent,
			},
		); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "error publishing booking event", zap.Error(err))
		}

		tempBooking := &models.Booking{SQLModel: common.SQLModel{Id: newBookingId}}
		tempBooking.Mask(false)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]string{
			"id": tempBooking.GetFakeId(),
		}))
	}
}
