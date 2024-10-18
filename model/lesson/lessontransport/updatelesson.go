// lesson/lessontransport/updatelesson.go

package lessontransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/lesson/lessonbiz"
	"video_server/model/lesson/lessonmodel"
	"video_server/model/lesson/lessonrepo"
	"video_server/model/lesson/lessonstore"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
)

func UpdateLessonHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lessonId := ctx.Param("id")

		var input lessonmodel.UpdateLesson
		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()

		courseStore := coursestore.NewSQLStore(db)
		lessonStore := lessonstore.NewSQLStore(db)
		videoStore := videostore.NewSQLStore(db)

		repo := lessonrepo.NewUpdateLessonRepo(lessonStore, courseStore, videoStore)
		biz := lessonbiz.NewUpdateLessonBiz(repo)

		if err := biz.UpdateLesson(ctx.Request.Context(), lessonId, &input); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
