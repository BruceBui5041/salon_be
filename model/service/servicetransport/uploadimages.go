package servicetransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/image/imagerepo"
	"salon_be/model/image/imagestore"
	"salon_be/model/m2mserviceversionimage/m2mserviceversionimagestore"
	"salon_be/model/service/servicebiz"
	"salon_be/model/service/servicemodel"
	"salon_be/model/service/servicerepo"
	"salon_be/model/serviceversion/serviceversionstore"

	"github.com/gin-gonic/gin"
)

func UploadImagesHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data servicemodel.UploadImages

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrNoPermission(errors.New("requester not found")))
		}
		data.UploadedBy = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		imageStore := imagestore.NewSQLStore(db)
		serviceVersionStore := serviceversionstore.NewSQLStore(db)
		m2mStore := m2mserviceversionimagestore.NewSQLStore(db)
		imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())

		repo := servicerepo.NewUploadImagesRepo(
			imageStore,
			imageRepo,
			serviceVersionStore,
			m2mStore,
			db,
		)
		biz := servicebiz.NewUploadImagesBiz(repo)

		if err := biz.UploadImages(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
