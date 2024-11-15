package usertransport

import (
	"net/http"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/hasher"
	"salon_be/component/tokenprovider/jwt"
	"salon_be/model/otp/otpbiz"
	"salon_be/model/otp/otpstore"
	"salon_be/model/user/userbiz"
	"salon_be/model/user/usermodel"
	"salon_be/model/user/userstore"
	"salon_be/utils"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(appCtx component.AppContext) func(*gin.Context) {
	return func(ctx *gin.Context) {
		db := appCtx.GetMainDBConnection()

		var data usermodel.CreateUser

		if err := ctx.ShouldBind(&data); err != nil {
			panic(err)
		}
		var userID string
		if err := db.Transaction(func(tx *gorm.DB) error {
			store := userstore.NewSQLStore(tx)
			otpStore := otpstore.NewSQLStore(tx)
			useStore := userstore.NewSQLStore(tx)

			md5 := hasher.NewMD5Hash()
			tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
			otpBiz := otpbiz.NewCreateOTPBiz(otpStore, useStore, appCtx.GetSMSClient())

			business := userbiz.NewRegisterBusiness(
				store,
				md5,
				tokenProvider,
				otpBiz,
			)

			account, user, err := business.RegisterUser(ctx.Request.Context(), &data, appconst.TokenExpiry)
			if err != nil {
				return err
			}

			userID = user.GetFakeId()

			utils.WriteServerJWTTokenCookie(ctx, account.Token.Token)

			ctx.JSON(http.StatusCreated, common.SimpleSuccessResponse(user.GetFakeId()))
			return nil
		}); err != nil {
			panic(err)
		}

		if err := watermill.PublishUserUpdated(
			ctx.Request.Context(),
			appCtx.GetLocalPubSub().GetUnblockPubSub(),
			&messagemodel.UserUpdatedMessage{UserId: userID},
		); err != nil {
			panic(err)
		}
	}
}
