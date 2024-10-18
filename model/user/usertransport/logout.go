package usertransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/user/userbiz"
	"video_server/utils"

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
