package servicetransport

import (
	"encoding/json"
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/model/image/imagerepo"
	"salon_be/model/image/imagestore"
	"salon_be/model/service/servicebiz"
	"salon_be/model/service/servicemodel"
	"salon_be/model/service/servicerepo"
	"salon_be/model/service/servicestore"
	"salon_be/model/serviceversion/serviceversionstore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func CreateServiceHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request servicemodel.CreateServiceRequest

		// Parse multipart form
		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
			logger.AppLogger.Error(ctx.Request.Context(), "failed to parse multipart form", zap.Error(err))
			panic(err)
		}

		// Bind the request struct (will get JSON string and images)
		if err := ctx.ShouldBind(&request); err != nil {
			logger.AppLogger.Error(ctx.Request.Context(), "failed to bind form data", zap.Error(err))
			panic(common.ErrInvalidRequest(err))
		}

		// Parse the JSON string into CreateService
		var serviceData servicemodel.CreateService
		if err := json.Unmarshal([]byte(request.JSON), &serviceData); err != nil {
			logger.AppLogger.Error(ctx.Request.Context(), "failed to unmarshal JSON data", zap.Error(err))
			panic(common.ErrInvalidRequest(err))
		}

		// Assign the uploaded images to the service version
		if serviceData.ServiceVersion != nil {
			serviceData.ServiceVersion.Images = request.Images
		}

		requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("requester not found")))
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			serviceStore := servicestore.NewSQLStore(tx)
			serviceVersionStore := serviceversionstore.NewSQLStore(tx)
			imageStore := imagestore.NewSQLStore(tx)

			imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())
			repo := servicerepo.NewCreateServiceRepo(serviceStore, serviceVersionStore, imageRepo)
			business := servicebiz.NewCreateServiceBiz(repo)

			serviceData.CreatorID = requester.GetUserId()
			if err := business.CreateNewService(ctx.Request.Context(), &serviceData); err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
