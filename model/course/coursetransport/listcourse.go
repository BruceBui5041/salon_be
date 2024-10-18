package coursetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursestore"

	"github.com/gin-gonic/gin"
)

func ListCourses(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		store := coursestore.NewSQLStore(db)
		biz := coursebiz.NewCourseBiz(store)

		result, err := biz.ListCourses(c.Request.Context())

		if err != nil {
			panic(err)
		}

		// for i := range result {
		// 	if i == len(result)-1 {
		// 	    paging.NextCursor = result[i].FakeId.String()
		// 	}
		// }

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
