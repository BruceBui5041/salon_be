package generictransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component/genericapi/genericbiz"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/component/genericapi/genericstore"

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
