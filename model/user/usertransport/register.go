package usertransport

import (
	"net/http"
	"video_server/appconst"
	"video_server/common"
	"video_server/component"
	"video_server/component/hasher"
	"video_server/component/tokenprovider/jwt"
	"video_server/model/user/userbiz"
	"video_server/model/user/usermodel"
	"video_server/model/user/userstore"
	"video_server/utils"
	"video_server/watermill"
	"video_server/watermill/messagemodel"

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
