package otptransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/otp/otpbiz"
	"salon_be/model/otp/otpmodel"
	"salon_be/model/otp/otprepo"
	"salon_be/model/otp/otpstore"
	"salon_be/model/user/userstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResendOTP(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input otpmodel.ResendOTPInput

		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("requester not found")))
		}
		input.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			otpStore := otpstore.NewSQLStore(tx)
			userStore := userstore.NewSQLStore(tx)
			repo := otprepo.NewResendRepo(otpStore, userStore)
			biz := otpbiz.NewResendOTPBiz(
				repo,
				appCtx.GetSMSClient(),
			)

			if err := biz.ResendOTP(c.Request.Context(), &input); err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
