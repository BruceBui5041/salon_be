package categorytransport

import (
	"errors"
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/category/categorybiz"
	"salon_be/model/category/categorymodel"
	"salon_be/model/category/categoryrepo"
	"salon_be/model/category/categorystore"

	"github.com/gin-gonic/gin"
)

func UpdateCategoryHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		var input categorymodel.UpdateCategory

		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := ctx.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}

		// Check if the requester is an admin (you might want to implement this check)
		if !requester.IsAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()

		store := categorystore.NewSQLStore(db)
		repo := categoryrepo.NewUpdateCategoryRepo(store, appCtx.GetS3Client())
		biz := categorybiz.NewUpdateCategoryBiz(repo)

		if err := biz.UpdateCategory(ctx.Request.Context(), uint32(id), &input); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
