package generictransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component/genericapi/genericbiz"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/component/genericapi/genericstore"

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
