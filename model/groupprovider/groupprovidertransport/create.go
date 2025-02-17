package groupprovidertransport

import (
	"encoding/json"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/model/groupprovider/groupproviderbiz"
	"salon_be/model/groupprovider/groupprovidermodel"
	"salon_be/model/groupprovider/groupproviderrepo"
	"salon_be/model/groupprovider/groupproviderstore"
	"salon_be/model/image/imagerepo"
	"salon_be/model/image/imagestore"
	"salon_be/model/service/servicestore"
	"salon_be/model/user/userstore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateGroupProvider(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request groupprovidermodel.GroupProviderCreateRequest

		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to parse multipart form", zap.Error(err))
			panic(err)
		}

		if err := c.ShouldBind(&request); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to bind form data", zap.Error(err))
			panic(common.ErrInvalidRequest(err))
		}

		var data groupprovidermodel.GroupProviderCreate
		if err := json.Unmarshal([]byte(request.JSON), &data); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to unmarshal JSON data", zap.Error(err))
			panic(common.ErrInvalidRequest(err))
		}

		data.Images = request.Images

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		data.RequesterID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			groupStore := groupproviderstore.NewSQLStore(tx)
			userStore := userstore.NewSQLStore(tx)
			serviceStore := servicestore.NewSQLStore(tx)
			imageStore := imagestore.NewSQLStore(tx)

			imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())
			repo := groupproviderrepo.NewCreateRepo(groupStore, userStore, serviceStore, imageRepo)
			biz := groupproviderbiz.NewCreateBiz(repo)

			if err := biz.CreateGroupProvider(c.Request.Context(), &data); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
