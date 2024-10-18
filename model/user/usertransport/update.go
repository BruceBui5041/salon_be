package usertransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/user/userbiz"
	"salon_be/model/user/usermodel"
	"salon_be/model/user/userrepo"
	"salon_be/model/user/userstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UpdateUser(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var input usermodel.UserUpdate

		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}

		// Check if the requester is the user being updated
		if requester.GetUserId() != uid.GetLocalID() {
			// Check if the requester is an admin or super admin
			if !requester.IsAdmin() && !requester.IsSuperAdmin() {
				panic(common.ErrNoPermission(nil))
			}
		}

		// Only allow role updates for admin or super admin
		if len(input.RoleIDs) > 0 && !requester.IsAdmin() && !requester.IsSuperAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			store := userstore.NewSQLStore(tx)
			repo := userrepo.NewUpdateUserRepo(store)
			biz := userbiz.NewUpdateUserBiz(repo)
			if err := biz.UpdateUser(c.Request.Context(), tx, uid.GetLocalID(), &input); err != nil {
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
