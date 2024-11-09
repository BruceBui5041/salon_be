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
