// File: lecturetransport/updatelecture.go

package lecturetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/lecture/lecturebiz"
	"video_server/model/lecture/lecturemodel"
	"video_server/model/lecture/lecturerepo"
	"video_server/model/lecture/lecturestore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateLectureHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lectureId := ctx.Param("id")

		var input lecturemodel.UpdateLecture
		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			lectureStore := lecturestore.NewSQLStore(tx)
			courseStore := coursestore.NewSQLStore(tx)

			repo := lecturerepo.NewUpdateLectureRepo(lectureStore, courseStore)
			biz := lecturebiz.NewUpdateLectureBiz(repo)

			if err := biz.UpdateLecture(ctx.Request.Context(), lectureId, &input); err != nil {
				return err
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}

	}
}
