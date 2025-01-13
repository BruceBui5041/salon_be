package usertransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/booking/bookingstore"
	"salon_be/model/payment/paymentstore"
	"salon_be/model/user/userbiz"
	"salon_be/model/user/usermodel"
	"salon_be/model/user/userrepo"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProviderEarnings(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var query usermodel.GetEarningsRequest
		if err := c.ShouldBind(&query); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		// Create date range for the specified month
		fromDate := time.Date(query.Year, time.Month(query.Month), 1, 0, 0, 0, 0, time.UTC)
		toDate := fromDate.AddDate(0, 1, 0).Add(-time.Second) // End of the month

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsProvider() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()
		bookingStore := bookingstore.NewSQLStore(db)
		paymentStore := paymentstore.NewSQLStore(db)

		repo := userrepo.NewProviderEarningsRepo(db, bookingStore, paymentStore)
		biz := userbiz.NewProviderEarningsBiz(repo)

		result, err := biz.GetProviderEarnings(
			c.Request.Context(),
			requester.GetUserId(),
			fromDate,
			toDate,
		)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
