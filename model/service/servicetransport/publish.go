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

func PublishServiceHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data servicemodel.PublishServiceRequest

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("user not authenticated")))
		}

		db := appCtx.GetMainDBConnection()
		if err := db.Transaction(func(tx *gorm.DB) error {
			serviceStore := servicestore.NewSQLStore(tx)
			serviceVersionStore := serviceversionstore.NewSQLStore(tx)

			repo := servicerepo.NewPublishServiceRepo(serviceStore, serviceVersionStore)
			business := servicebiz.NewPublishServiceBiz(repo)

			if err := business.PublishService(c.Request.Context(), requester, &data); err != nil {
				panic(err)
			}

			c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
