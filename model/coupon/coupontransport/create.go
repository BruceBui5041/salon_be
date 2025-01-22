package coupontransport

import (
	"encoding/json"
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	models "salon_be/model"
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

func CreateCouponHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request couponmodel.CreateCouponRequest

		// Parse multipart form
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to parse multipart form", zap.Error(err))
			panic(err)
		}

		// Bind the request struct
		if err := c.ShouldBind(&request); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to bind form data", zap.Error(err))
			panic(couponerror.ErrCouponInvalid(err))
		}

		// Parse the JSON string into CreateCoupon
		var data couponmodel.CreateCoupon
		if err := json.Unmarshal([]byte(request.JSON), &data); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "failed to unmarshal JSON data", zap.Error(err))
			panic(couponerror.ErrCouponInvalid(err))
		}

		if data.CatergoryStrId != nil {
			cateUID, err := common.FromBase58(*data.CatergoryStrId)
			if err != nil {
				panic(couponerror.ErrCouponInvalid(err))
			}
			cateId := cateUID.GetLocalID()
			data.CategoryID = &cateId
		}

		// Assign the uploaded image
		data.Image = request.Image

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			err := errors.New("invalid user type")
			logger.AppLogger.Error(c.Request.Context(), "Invalid user type in create coupon",
				zap.Error(err),
				zap.Any("requester", requester))
			panic(common.ErrNoPermission(err))
		}

		if !requester.IsAdmin() && !requester.IsSuperAdmin() {
			err := errors.New("user must be admin or super admin")
			logger.AppLogger.Error(c.Request.Context(), "Insufficient permissions for create coupon",
				zap.Error(err),
				zap.Uint32("user_id", requester.GetUserId()),
				zap.Any("user_role", requester.GetRoles(c.Request.Context())))
			panic(common.ErrNoPermission(err))
		}

		data.CreatorID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		var couponId uint32
		if err := db.Transaction(func(tx *gorm.DB) error {
			store := couponstore.NewSQLStore(tx)
			imageStore := imagestore.NewSQLStore(tx)
			imageRepo := imagerepo.NewCreateImageRepo(imageStore, appCtx.GetS3Client())
			repo := couponrepo.NewCreateCouponRepo(store, imageRepo)
			biz := couponbiz.NewCreateCouponBiz(repo)

			id, err := biz.CreateCoupon(c.Request.Context(), &data)
			if err != nil {
				logger.AppLogger.Error(c.Request.Context(), "Failed to create coupon in handler",
					zap.Error(err),
					zap.Any("coupon_data", data))
				return err
			}
			couponId = id
			return nil
		}); err != nil {
			panic(err)
		}

		tempCoupon := &models.Coupon{SQLModel: common.SQLModel{Id: couponId}}
		tempCoupon.Mask(false)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]string{
			"id": tempCoupon.GetFakeId(),
		}))
	}
}
