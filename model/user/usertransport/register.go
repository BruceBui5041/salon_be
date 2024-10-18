package usertransport

import (
	"net/http"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/hasher"
	"salon_be/component/tokenprovider/jwt"
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

		if err := db.Transaction(func(tx *gorm.DB) error {
			store := userstore.NewSQLStore(tx)
			md5 := hasher.NewMD5Hash()
			tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

			business := userbiz.NewRegisterBusiness(
				store,
				md5,
				tokenProvider,
			)

			account, user, err := business.RegisterUser(ctx.Request.Context(), &data, appconst.TokenExpiry)
			if err != nil {
				return err
			}

			if err := watermill.PublishUserUpdated(
				ctx.Request.Context(),
				appCtx.GetLocalPubSub().GetUnblockPubSub(),
				&messagemodel.UserUpdatedMessage{UserId: user.GetFakeId()},
			); err != nil {
				return common.ErrInternal(err)
			}

			utils.WriteServerJWTTokenCookie(ctx, account.Token)

			data.Mask(false)

			ctx.JSON(http.StatusCreated, common.SimpleSuccessResponse(data.FakeId.String()))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
