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
)

func UpdatePermissionHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			panic(common.ErrInvalidRequest(errors.New("id missing")))
		}

		var input permissionmodel.UpdatePermission
		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}

		if !requester.IsSuperAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()
		store := permissionstore.NewSQLStore(db)
		repo := permissionrepo.NewUpdatePermissionRepo(store)
		biz := permissionbiz.NewUpdatePermissionBiz(repo)

		if err := biz.UpdatePermission(ctx, id, &input); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
