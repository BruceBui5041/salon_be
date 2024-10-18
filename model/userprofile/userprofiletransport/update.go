package userprofiletransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/userprofile/userprofilebiz"
	"salon_be/model/userprofile/userprofilemodel"
	"salon_be/model/userprofile/userprofilerepo"
	"salon_be/model/userprofile/userprofilestore"

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
