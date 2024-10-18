package lessontransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/lesson/lessonbiz"
	"video_server/model/lesson/lessonstore"

	"github.com/gin-gonic/gin"
)

func GetLessonHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		db := appCtx.GetMainDBConnection()
		store := lessonstore.NewSQLStore(db)
		biz := lessonbiz.NewGetLessonBiz(store)

		lesson, err := biz.GetLessonByID(c.Request.Context(), id)
		if err != nil {
			panic(err)
		}

		lesson.Mask(false)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(lesson))
	}
}
