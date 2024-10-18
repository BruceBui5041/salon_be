package coursetransport

import (
	"errors"
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/category/categorystore"
	"video_server/model/course/coursebiz"
	"video_server/model/course/coursemodel"
	"video_server/model/course/courserepo"
	"video_server/model/course/coursestore"
	"video_server/model/video/videostore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateCourseHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		var input coursemodel.UpdateCourse
		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInternal(err))
		}

		requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInternal(errors.New("cannot find requester")))
		}

		requester.Mask(false)

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			input.UploadedBy = requester.GetFakeId()
			input.Id = id
			input.Mask(false)

			svc := appCtx.GetS3Client()

			categoryStore := categorystore.NewSQLStore(tx)
			courseStore := coursestore.NewSQLStore(tx)
			videoStore := videostore.NewSQLStore(tx)
			repo := courserepo.NewUpdateCourseRepo(courseStore, categoryStore, videoStore, svc)
			courseBusiness := coursebiz.NewUpdateCourseBiz(repo)

			err = courseBusiness.UpdateCourse(ctx.Request.Context(), uint32(id), &input)
			if err != nil {
				return err
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("ok"))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
