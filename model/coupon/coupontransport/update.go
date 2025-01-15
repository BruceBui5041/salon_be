package coupontransport

import (
	"encoding/json"
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/model/coupon/couponbiz"
	"salon_be/model/coupon/couponerror"
	"salon_be/model/coupon/couponmodel"
	"salon_be/model/coupon/couponrepo"
	"salon_be/model/coupon/couponstore"
	"salon_be/model/image/imagerepo"
	"salon_be/model/image/imagestore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UpdateCouponHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var request couponmodel.UpdateCouponRequest

		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to parse multipart form", zap.Error(err))
			panic(err)
		}

		if err := c.ShouldBind(&request); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to bind form data", zap.Error(err))
			panic(couponerror.ErrCouponInvalid(err))
		}

		var data couponmodel.UpdateCoupon
		if err := json.Unmarshal([]byte(request.JSON), &data); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to unmarshal JSON data", zap.Error(err))
			panic(couponerror.ErrCouponInvalid(err))
		}

		data.Image = request.Image

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsAdmin() && !requester.IsSuperAdmin() {
			err := errors.New("user must be admin or super admin")
			logger.AppLogger.Error(c.Request.Context(), "Insufficient permissions for update coupon",
				zap.Error(err),
				zap.Uint32("user_id", requester.GetUserId()),
				zap.Any("user_role", requester.GetRoles(c.Request.Context())))
			panic(common.ErrNoPermission(err))
		}

		db := appCtx.GetMainDBConnection()

		data.CreatorID = requester.GetUserId()
		if err := db.Transaction(func(tx *gorm.DB) error {
			store := couponstore.NewSQLStore(tx)
			imageStore := imagestore.NewSQLStore(tx)
			imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())
			repo := couponrepo.NewUpdateCouponRepo(store, imageRepo)
			biz := couponbiz.NewUpdateCouponBiz(repo)

			if err := biz.UpdateCoupon(c.Request.Context(), id, &data); err != nil {
				logger.AppLogger.Error(c.Request.Context(), "Failed to update coupon in handler",
					zap.Error(err),
					zap.String("coupon_id", id),
					zap.Any("coupon_data", data))
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
