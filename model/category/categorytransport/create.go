package categorytransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/category/categorybiz"
	"salon_be/model/category/categorymodel"
	"salon_be/model/category/categoryrepo"
	"salon_be/model/category/categorystore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCategoryHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input categorymodel.CreateCategory

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		// Check if the requester is an admin (you might want to implement this check)
		if !requester.IsAdmin() {
			panic(common.ErrNoPermission(nil))
		}

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			categoryStore := categorystore.NewSQLStore(tx)
			repo := categoryrepo.NewCreateCategoryRepo(categoryStore)
			categoryBusiness := categorybiz.NewCreateCategoryBiz(repo)

			if err := categoryBusiness.CreateNewCategory(ctx.Request.Context(), &input); err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(input.GetFakeId()))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
