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

func CreateLectureHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input lecturemodel.CreateLecture

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			courseStore := coursestore.NewSQLStore(tx)
			lectureStore := lecturestore.NewSQLStore(tx)

			repo := lecturerepo.NewCreateLectureRepo(lectureStore, courseStore)
			lectureBusiness := lecturebiz.NewCreateLectureBiz(repo)

			lecture, err := lectureBusiness.CreateNewLecture(ctx.Request.Context(), &input)

			if err != nil {
				return err
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(lecture))
			return nil
		}); err != nil {
			panic(err)
		}

	}
}
