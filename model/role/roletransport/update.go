package roletransport

import (
	"errors"
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/role/rolebiz"
	"video_server/model/role/rolemodel"
	"video_server/model/role/rolerepo"
	"video_server/model/role/rolestore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateRoleHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		var input rolemodel.UpdateRole
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

		if err := db.Transaction(func(tx *gorm.DB) error {
			roleStore := rolestore.NewSQLStore(tx)
			repo := rolerepo.NewUpdateRoleRepo(roleStore)
			roleBusiness := rolebiz.NewUpdateRoleBiz(repo)

			if err := roleBusiness.UpdateRole(ctx, uint32(id), &input); err != nil {
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
