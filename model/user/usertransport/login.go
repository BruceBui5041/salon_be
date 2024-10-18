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

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginUserData usermodel.UserLogin

		if err := ctx.ShouldBind(&loginUserData); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

			md5 := hasher.NewMD5Hash()

			userStore := userstore.NewSQLStore(tx)
			loginbiz := userbiz.NewLoginBusiness(
				userStore,
				tokenProvider,
				md5,
				appconst.TokenExpiry, // 7 days
				appCtx.GetLocalPubSub().GetBlockPubSub(),
			)

			loginRes, err := loginbiz.Login(ctx.Request.Context(), &loginUserData)
			if err != nil {
				return err
			}

			utils.WriteServerJWTTokenCookie(ctx, loginRes.Token.Token)
			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(loginRes))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
