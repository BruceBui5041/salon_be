package roletransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/role/rolebiz"
	"video_server/model/role/rolerepo"
	"video_server/model/role/rolestore"

	"github.com/gin-gonic/gin"
)

func SoftDeleteRoleHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		if !requester.IsSuperAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()

		store := rolestore.NewSQLStore(db)
		repo := rolerepo.NewDeleteRoleRepo(store)
		biz := rolebiz.NewDeleteRoleBiz(repo)

		if err := biz.SoftDeleteRole(ctx, uint32(id)); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
