package paymenttransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/enrollment/enrollmentrepo"
	"video_server/model/enrollment/enrollmentstore"
	"video_server/model/payment/paymentbiz"
	"video_server/model/payment/paymentmodel"
	"video_server/model/payment/paymentrepo"
	"video_server/model/payment/paymentstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePaymentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input paymentmodel.CreatePayment

		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		input.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			paymentStore := paymentstore.NewSQLStore(tx)
			courseStore := coursestore.NewSQLStore(tx)
			enrollmentStore := enrollmentstore.NewSQLStore(tx)
			paymentRepo := paymentrepo.NewCreatePaymentRepo(paymentStore)
			enrollmentRepo := enrollmentrepo.NewCreateEnrollmentRepo(
				enrollmentStore, courseStore,
				appCtx.GetLocalPubSub().GetUnblockPubSub(),
			)
			biz := paymentbiz.NewCreatePaymentBiz(paymentRepo, enrollmentRepo)

			payment, err := biz.CreateNewPayment(c.Request.Context(), &input)
			if err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(payment))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
