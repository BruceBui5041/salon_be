package permissiontransport

import (
	"errors"
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/permission/permissionbiz"
	"video_server/model/permission/permissionmodel"
	"video_server/model/permission/permissionrepo"
	"video_server/model/permission/permissionstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreatePermissionHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input permissionmodel.CreatePermission

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}

		if !requester.IsSuperAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			permissionStore := permissionstore.NewSQLStore(tx)
			repo := permissionrepo.NewCreatePermissionRepo(permissionStore)
			permissionBusiness := permissionbiz.NewCreatePermissionBiz(repo)

			if err := permissionBusiness.CreateNewPermission(ctx, &input); err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
