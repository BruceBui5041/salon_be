package generictransport

import (
	"net/http"
	"video_server/common"
	"video_server/component/genericapi/genericbiz"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/genericapi/genericstore"

	"github.com/gin-gonic/gin"
)

func (gt *GenericTransport) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input genericmodel.SearchModelRequest
		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInternal(err))
		}

		db := gt.AppContext.GetMainDBConnection()
		store := genericstore.NewGenericStore(db)
		biz := genericbiz.NewGenericBiz(store)

		var resp []interface{}
		if err := biz.Search(c.Request.Context(), input, &resp); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
	}
}
