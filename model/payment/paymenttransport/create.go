package paymenttransport

import (
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/payment/paymentmodel"

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
			// paymentStore := paymentstore.NewSQLStore(tx)
			// courseStore := coursestore.NewSQLStore(tx)
			// enrollmentStore := enrollmentstore.NewSQLStore(tx)
			// paymentRepo := paymentrepo.NewCreatePaymentRepo(paymentStore)
			// enrollmentRepo := enrollmentrepo.NewCreateEnrollmentRepo(
			// 	enrollmentStore, courseStore,
			// 	appCtx.GetLocalPubSub().GetUnblockPubSub(),
			// )
			// biz := paymentbiz.NewCreatePaymentBiz(paymentRepo, enrollmentRepo)

			// payment, err := biz.CreateNewPayment(c.Request.Context(), &input)
			// if err != nil {
			// 	panic(err)
			// }

			// c.JSON(http.StatusOK, common.SimpleSuccessResponse(payment))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
