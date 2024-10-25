// service/servicetransport/create.go
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

func CreateServiceHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input servicemodel.CreateService

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("requester not found")))
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			serviceStore := servicestore.NewSQLStore(tx)
			serviceVersionStore := serviceversionstore.NewSQLStore(tx)
			repo := servicerepo.NewCreateServiceRepo(serviceStore, serviceVersionStore)
			business := servicebiz.NewCreateServiceBiz(repo)

			input.CreatorID = requester.GetUserId()
			if err := business.CreateNewService(ctx.Request.Context(), &input); err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
