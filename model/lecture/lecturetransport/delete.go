package lecturetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/lecture/lecturebiz"
	"video_server/model/lecture/lecturerepo"
	"video_server/model/lecture/lecturestore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteLectureHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lectureId := ctx.Param("id")

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			lectureStore := lecturestore.NewSQLStore(tx)
			courseStore := coursestore.NewSQLStore(tx)

			repo := lecturerepo.NewDeleteLectureRepo(lectureStore, courseStore)
			biz := lecturebiz.NewDeleteLectureBiz(repo)

			if err := biz.DeleteLecture(ctx.Request.Context(), lectureId); err != nil {
				return err
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}

	}
}
