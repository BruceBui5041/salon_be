package enrollmenttransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/course/coursestore"
	"video_server/model/enrollment/enrollmentbiz"
	"video_server/model/enrollment/enrollmentmodel"
	"video_server/model/enrollment/enrollmentrepo"
	"video_server/model/enrollment/enrollmentstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateEnrollmentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input enrollmentmodel.CreateEnrollment

		if err := c.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			enrollmentStore := enrollmentstore.NewSQLStore(tx)
			courseStore := coursestore.NewSQLStore(tx)
			enrollmentRepo := enrollmentrepo.NewCreateEnrollmentRepo(
				enrollmentStore,
				courseStore,
				appCtx.GetLocalPubSub().GetUnblockPubSub(),
			)
			biz := enrollmentbiz.NewCreateEnrollmentBiz(enrollmentRepo)

			err := biz.CreateNewEnrollment(c.Request.Context(), &input)
			if err != nil {
				panic(err)
			}

			requester.Mask(false)

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(common.ErrInternal(err))
		}
	}
}
