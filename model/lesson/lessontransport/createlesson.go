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

func CreateLessonHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input lessonmodel.CreateLesson

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()

		courseStore := coursestore.NewSQLStore(db)
		videoStore := videostore.NewSQLStore(db)
		lessonStore := lessonstore.NewSQLStore(db)

		repo := lessonrepo.NewCreateLessonRepo(lessonStore, courseStore, videoStore)
		lessonBusiness := lessonbiz.NewCreateLessonBiz(repo)

		lesson, err := lessonBusiness.CreateNewLesson(ctx.Request.Context(), &input)

		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(lesson))
	}
}
