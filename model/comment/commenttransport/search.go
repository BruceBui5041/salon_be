package commenttransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/model/comment/commentbiz"
	"salon_be/model/comment/commentstore"

	"github.com/gin-gonic/gin"
)

func SearchComment(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input genericmodel.SearchModelRequest
		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInternal(err))
		}

		db := appCtx.GetMainDBConnection()
		store := commentstore.NewSQLStore(db)
		biz := commentbiz.NewSearchCommentBiz(store)

		result, err := biz.SearchComments(c.Request.Context(), input)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
