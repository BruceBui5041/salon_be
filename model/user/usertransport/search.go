package usertransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/model/user/userbiz"
	"salon_be/model/user/userstore"

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
