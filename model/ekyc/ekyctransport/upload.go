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

func UploadImage(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var request ekycmodel.UploadRequest
		if err := c.ShouldBind(&request); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsUser() {
			panic(common.ErrNoPermission(errors.New("only users can upload KYC images")))
		}

		db := appCtx.GetMainDBConnection()
		var result ekycmodel.KYCImageUploadRes

		if err := db.Transaction(func(tx *gorm.DB) error {
			imageStore := imagestore.NewSQLStore(tx)
			imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())

			kycStore := ekycstore.NewSQLStore(tx)
			repo := ekycrepo.NewKYCImageUploadRepo(
				kycStore,
				appCtx.GetEKYCClient(),
				imageRepo,
			)
			biz := ekycbiz.NewUploadBiz(repo)

			res, err := biz.UploadImage(c.Request.Context(), requester.GetUserId(), &request)
			if err != nil {
				return err
			}

			result = *res
			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
