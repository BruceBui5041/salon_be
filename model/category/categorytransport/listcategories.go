package categorytransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/category/categorybiz"
	"salon_be/model/category/categorystore"

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
			"Services",
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
