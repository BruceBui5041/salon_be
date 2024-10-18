package lessontransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/lesson/lessonbiz"
	"video_server/model/lesson/lessonrepo"
	"video_server/model/lesson/lessonstore"

	"github.com/gin-gonic/gin"
)

func DeleteLessonHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lessonId := ctx.Param("id")

		db := appCtx.GetMainDBConnection()

		lessonStore := lessonstore.NewSQLStore(db)
		courseStore := coursestore.NewSQLStore(db)

		repo := lessonrepo.NewDeleteLessonRepo(lessonStore, courseStore)
		biz := lessonbiz.NewDeleteLessonBiz(repo)

		if err := biz.DeleteLesson(ctx.Request.Context(), lessonId); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
