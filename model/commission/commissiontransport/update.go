package commissiontransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	commissionmodel "salon_be/model/commission/comissionmodel"
	"salon_be/model/commission/commissionbiz"
	"salon_be/model/commission/commissionrepo"
	"salon_be/model/commission/commissionstore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func UpdateCommissionHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		var data commissionmodel.UpdateCommission
		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		if !requester.IsAdmin() && !requester.IsSuperAdmin() {
			err := errors.New("user must be admin or super admin")
			logger.AppLogger.Error(
				c.Request.Context(),
				"Insufficient permissions for update commission",
				zap.Error(err),
				zap.Uint32("user_id", requester.GetUserId()),
			)
			panic(common.ErrNoPermission(err))
		}

		roleUID, err := common.FromBase58(data.RoleIDStr)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		data.RoleID = roleUID.GetLocalID()

		data.UpdaterID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			store := commissionstore.NewSQLStore(tx)
			repo := commissionrepo.NewUpdateCommissionRepo(store)
			biz := commissionbiz.NewUpdateCommissionBiz(repo)

			if err := biz.UpdateCommission(c.Request.Context(), id, &data); err != nil {
				logger.AppLogger.Error(c.Request.Context(), "Failed to update commission in handler",
					zap.Error(err),
					zap.Uint32("id", id),
					zap.Any("commission_data", data))
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
