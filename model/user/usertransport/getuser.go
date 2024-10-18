package usertransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/user/userbiz"
	"video_server/model/user/userrepo"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func GetUser(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester := c.MustGet(common.CurrentUser).(common.Requester)
		store := userstore.NewSQLStore(appCtx.GetMainDBConnection())
		repo := userrepo.NewGetUserRepo(store)
		biz := userbiz.NewGetUserBiz(repo)

		user, err := biz.GetUserById(c.Request.Context(), requester.GetUserId())
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(user))
	}
}
