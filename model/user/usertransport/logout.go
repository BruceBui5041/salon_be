package usertransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/user/userbiz"
	"salon_be/utils"

	"github.com/gin-gonic/gin"
)

func Logout(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logoutBiz := userbiz.NewLogoutBusiness(appCtx.GetAppCache())

		err := logoutBiz.Logout(ctx.Request.Context())
		if err != nil {
			panic(err)
		}

		// Clear the JWT cookie
		utils.ClearServerJWTTokenCookie(ctx)

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
