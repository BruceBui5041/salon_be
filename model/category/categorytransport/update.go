package categorytransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/category/categorybiz"
	"salon_be/model/category/categorymodel"
	"salon_be/model/category/categoryrepo"
	"salon_be/model/category/categorystore"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateCategoryHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var input categorymodel.UpdateCategory

		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		// Check if the requester is an admin (you might want to implement this check)
		if !requester.IsAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()

		store := categorystore.NewSQLStore(db)
		repo := categoryrepo.NewUpdateCategoryRepo(store)
		biz := categorybiz.NewUpdateCategoryBiz(repo)

		if err := biz.UpdateCategory(ctx.Request.Context(), uint32(id), &input); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
