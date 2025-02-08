package ekyctransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/ekyc/ekycbiz"
	"salon_be/model/ekyc/ekycmodel"
	"salon_be/model/ekyc/ekycrepo"
	"salon_be/model/ekyc/ekycstore"
	"salon_be/model/image/imagerepo"
	"salon_be/model/image/imagestore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateKYCProfile(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var request ekycmodel.CreateKYCProfileRequest
		if err := c.ShouldBind(&request); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsUser() {
			panic(common.ErrNoPermission(errors.New("only users can create KYC profiles")))
		}

		request.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			// Initialize stores with transaction
			imageStore := imagestore.NewSQLStore(tx)
			imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())

			store := ekycstore.NewSQLStore(tx)

			// Initialize upload repo
			uploadRepo := ekycrepo.NewKYCImageUploadRepo(
				store,
				appCtx.GetEKYCClient(),
				imageRepo,
			)

			// Initialize create repo with upload repo
			createRepo := ekycrepo.NewCreateKYCRepo(
				store,
				appCtx.GetEKYCClient(),
				uploadRepo,
			)

			// Initialize business logic
			business := ekycbiz.NewCreateKYCBiz(createRepo)

			if err := business.CreateKYCProfile(c.Request.Context(), &request); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
