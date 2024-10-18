package coursetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursestore"

	"github.com/gin-gonic/gin"
)

func GetCourseByID(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		db := appCtx.GetMainDBConnection()
		store := coursestore.NewSQLStore(db)
		biz := coursebiz.NewGetCourseByIDBiz(store)

		result, err := biz.GetCourseByID(c.Request.Context(), int(id))
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
