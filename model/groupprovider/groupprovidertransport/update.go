package groupprovidertransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/groupprovider/groupproviderbiz"
	"salon_be/model/groupprovider/groupprovidermodel"
	"salon_be/model/groupprovider/groupproviderrepo"
	"salon_be/model/groupprovider/groupproviderstore"
	"salon_be/model/user/userstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateGroupProvider(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var data groupprovidermodel.GroupProviderUpdate
		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		data.RequesterID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			groupStore := groupproviderstore.NewSQLStore(tx)
			userStore := userstore.NewSQLStore(tx)

			repo := groupproviderrepo.NewUpdateRepo(groupStore, userStore)
			biz := groupproviderbiz.NewUpdateBiz(repo)

			if err := biz.UpdateGroupProvider(c.Request.Context(), uid.GetLocalID(), &data); err != nil {
				return err
			}

			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
