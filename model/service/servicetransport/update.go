package servicetransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/service/servicebiz"
	"salon_be/model/service/servicemodel"
	"salon_be/model/service/servicerepo"
	"salon_be/model/service/servicestore"
	"salon_be/model/serviceversion/serviceversionstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateServiceHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input servicemodel.UpdateService

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		// requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		// if !ok {
		// 	panic(common.ErrInvalidRequest(errors.New("requester not found")))
		// }

		serviceUID, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(errors.New("invalid category ID")))
		}

		input.ServiceID = serviceUID.GetLocalID()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			serviceStore := servicestore.NewSQLStore(tx)
			serviceVersionStore := serviceversionstore.NewSQLStore(tx)
			repo := servicerepo.NewUpdateServiceRepo(serviceStore, serviceVersionStore)
			business := servicebiz.NewUpdateServiceBiz(repo)

			if err := business.UpdateService(ctx.Request.Context(), &input); err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
