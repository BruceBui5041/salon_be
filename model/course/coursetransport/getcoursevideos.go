package coursetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursestore"

	"github.com/gin-gonic/gin"
)

func GetCourseVideos(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		db := appCtx.GetMainDBConnection()
		store := coursestore.NewSQLStore(db)
		biz := coursebiz.NewGetCourseVideosBiz(store)

		result, err := biz.GetCourseVideos(c.Request.Context(), int(id))
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
