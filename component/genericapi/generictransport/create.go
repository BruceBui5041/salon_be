package generictransport

import (
	"net/http"
	"video_server/common"
	"video_server/component/genericapi/genericbiz"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/genericapi/genericstore"

	"github.com/gin-gonic/gin"
)

func (gt *GenericTransport) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input genericmodel.CreateRequest
		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := gt.AppContext.GetMainDBConnection()
		store := genericstore.NewGenericStore(db)
		biz := genericbiz.NewGenericBiz(store)

		result, err := biz.Create(c.Request.Context(), input)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
