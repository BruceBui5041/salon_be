package coursetransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/category/categorystore"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursemodel"
	"video_server/model/course/courserepo"
	"video_server/model/course/coursestore"
	"video_server/model/user/userstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCourseHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input coursemodel.CreateCourse

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			svc := appCtx.GetS3Client()

			categoryStore := categorystore.NewSQLStore(tx)
			coursestore := coursestore.NewSQLStore(tx)
			userStore := userstore.NewSQLStore(tx)
			repo := courserepo.NewCreateCourseRepo(
				coursestore,
				categoryStore,
				userStore,
				svc,
			)
			coursebusiness := coursebiz.NewCreateCourseBiz(repo)

			input.CreatorID = requester.GetUserId()
			course, err := coursebusiness.CreateNewCourse(ctx.Request.Context(), &input)
			if err != nil {
				return err
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(course))
			return nil

		}); err != nil {
			panic(err)
		}

	}
}
