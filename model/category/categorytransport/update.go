package categorytransport

import (
	"net/http"
	"strconv"
	"video_server/common"
	"video_server/component"
	"video_server/model/category/categorybiz"
	"video_server/model/category/categorymodel"
	"video_server/model/category/categoryrepo"
	"video_server/model/category/categorystore"

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
