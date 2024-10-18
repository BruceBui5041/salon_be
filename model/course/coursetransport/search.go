package coursetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/component/genericapi/genericbiz"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/genericapi/genericstore"
	"video_server/component/genericapi/generictransport"

	"github.com/gin-gonic/gin"
)

type courseTransport struct {
	generictransport.GenericTransport
}

func NewCourseTransport(appCtx component.AppContext) *courseTransport {
	return &courseTransport{
		GenericTransport: generictransport.GenericTransport{
			AppContext: appCtx,
		},
	}
}

func (ct *courseTransport) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input genericmodel.SearchModelRequest
		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInternal(err))
		}

		db := ct.AppContext.GetMainDBConnection()
		store := genericstore.NewGenericStore(db)
		biz := genericbiz.NewGenericBiz(store)

		var resp []interface{}
		if err := biz.Search(c.Request.Context(), input, &resp); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(resp))
	}
}
