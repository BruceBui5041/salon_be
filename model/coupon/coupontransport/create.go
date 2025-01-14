package coupontransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	models "salon_be/model"
	"salon_be/model/coupon/couponbiz"
	"salon_be/model/coupon/couponmodel"
	"salon_be/model/coupon/couponrepo"
	"salon_be/model/coupon/couponstore"

	"github.com/gin-gonic/gin"
)

func CreateCouponHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data couponmodel.CreateCoupon

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		user, ok := requester.(*models.User)
		if !ok {
			panic(common.ErrNoPermission(errors.New("invalid user type")))
		}

		if !user.IsAdmin() && !user.IsSuperAdmin() {
			panic(common.ErrNoPermission(errors.New("user must be admin or super admin")))
		}

		data.CreatorID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		store := couponstore.NewSQLStore(db)
		repo := couponrepo.NewCreateCouponRepo(store)
		biz := couponbiz.NewCreateCouponBiz(repo)

		if err := biz.CreateCoupon(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
