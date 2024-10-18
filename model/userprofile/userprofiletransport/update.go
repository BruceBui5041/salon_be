package userprofiletransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/userprofile/userprofilebiz"
	"video_server/model/userprofile/userprofilemodel"
	"video_server/model/userprofile/userprofilerepo"
	"video_server/model/userprofile/userprofilestore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateProfileHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input userprofilemodel.UpdateProfileModel
		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			store := userprofilestore.NewSQLStore(db)
			repo := userprofilerepo.NewUpdateProfileRepo(store, appCtx.GetS3Client())
			biz := userprofilebiz.NewUpdateProfileBiz(repo)

			if err := biz.UpdateProfile(
				ctx.Request.Context(),
				appCtx.GetLocalPubSub().GetBlockPubSub(),
				&input,
			); err != nil {
				return nil
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}

	}
}
