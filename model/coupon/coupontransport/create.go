package coupontransport

import (
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

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateCouponHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data couponmodel.CreateCoupon

		if err := c.ShouldBind(&data); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "Invalid request body for create coupon",
				zap.Error(err))
			panic(couponerror.ErrCouponInvalid(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		user, ok := requester.(*models.User)
		if !ok {
			err := errors.New("invalid user type")
			logger.AppLogger.Error(c.Request.Context(), "Invalid user type in create coupon",
				zap.Error(err),
				zap.Any("requester", requester))
			panic(common.ErrNoPermission(err))
		}

		if !user.IsAdmin() && !user.IsSuperAdmin() {
			err := errors.New("user must be admin or super admin")
			logger.AppLogger.Error(c.Request.Context(), "Insufficient permissions for create coupon",
				zap.Error(err),
				zap.Uint32("user_id", user.Id),
				zap.Any("user_role", user.GetRoles(c.Request.Context())))
			panic(common.ErrNoPermission(err))
		}

		data.CreatorID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		store := couponstore.NewSQLStore(db)
		repo := couponrepo.NewCreateCouponRepo(store)
		biz := couponbiz.NewCreateCouponBiz(repo)

		if err := biz.CreateCoupon(c.Request.Context(), &data); err != nil {
			logger.AppLogger.Error(c.Request.Context(), "Failed to create coupon in handler",
				zap.Error(err),
				zap.Any("coupon_data", data))
			panic(err)
		}

		logger.AppLogger.Info(c.Request.Context(), "Coupon created successfully",
			zap.String("code", data.Code),
			zap.Uint32("creator_id", data.CreatorID))

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
