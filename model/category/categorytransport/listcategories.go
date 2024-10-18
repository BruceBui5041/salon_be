package categorytransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/category/categorybiz"
	"video_server/model/category/categorystore"

	"github.com/gin-gonic/gin"
)

func ListCategories(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		store := categorystore.NewSQLStore(db)
		biz := categorybiz.NewCategoryBiz(store)

		result, err := biz.ListCategories(
			c.Request.Context(),
			map[string]interface{}{},
			"Courses",
		)

		if err != nil {
			panic(err)
		}

		// for i := range result {
		// result[i].Mask(false)

		// if i == len(result)-1 {
		// 	paging.NextCursor = result[i].FakeId.String()
		// }
		// }

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
