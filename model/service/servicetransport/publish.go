package servicetransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/service/servicebiz"
	"salon_be/model/service/servicemodel"
	"salon_be/model/service/servicerepo"
	"salon_be/model/service/servicestore"
	"salon_be/model/serviceversion/serviceversionstore"

	"github.com/gin-gonic/gin"
)

func PublishServiceHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data servicemodel.PublishServiceRequest

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()
		serviceStore := servicestore.NewSQLStore(db)
		serviceVersionStore := serviceversionstore.NewSQLStore(db)

		repo := servicerepo.NewPublishServiceRepo(serviceStore, serviceVersionStore)
		business := servicebiz.NewPublishServiceBiz(repo)

		if err := business.PublishService(c.Request.Context(), requester, &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
