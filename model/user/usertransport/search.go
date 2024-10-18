package usertransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/component/genericapi/genericmodel"
	"video_server/model/user/userbiz"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
)

func SearchUser(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input genericmodel.SearchModelRequest
		if err := c.ShouldBind(&input); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()
		store := userstore.NewSQLStore(db)
		biz := userbiz.NewSearchUserBiz(store)

		result, err := biz.SearchUsers(c.Request.Context(), input)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
