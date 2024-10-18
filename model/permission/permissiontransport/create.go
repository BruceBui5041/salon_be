package permissiontransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/permission/permissionbiz"
	"salon_be/model/permission/permissionmodel"
	"salon_be/model/permission/permissionrepo"
	"salon_be/model/permission/permissionstore"

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
