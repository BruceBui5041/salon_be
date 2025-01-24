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

func CreateCommissionHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data commissionmodel.CreateCommission

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		if !requester.IsAdmin() && !requester.IsSuperAdmin() {
			err := errors.New("user must be admin or super admin")
			logger.AppLogger.Error(c.Request.Context(), "Insufficient permissions for create commission",
				zap.Error(err),
				zap.Uint32("user_id", requester.GetUserId()))
			panic(common.ErrNoPermission(err))
		}

		data.CreatorID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		var commissionId uint32
		if err := db.Transaction(func(tx *gorm.DB) error {
			store := commissionstore.NewSQLStore(tx)
			repo := commissionrepo.NewCreateCommissionRepo(store)
			biz := commissionbiz.NewCreateCommissionBiz(repo)

			id, err := biz.CreateCommission(c.Request.Context(), &data)
			if err != nil {
				logger.AppLogger.Error(c.Request.Context(), "Failed to create commission in handler",
					zap.Error(err),
					zap.Any("commission_data", data))
				return err
			}
			commissionId = id
			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]interface{}{
			"id": commissionId,
		}))
	}
}
