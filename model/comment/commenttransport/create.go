package commenttransport

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"salon_be/model/comment/commentbiz"
	"salon_be/model/comment/commentmodel"
	"salon_be/model/comment/commentrepo"
	"salon_be/model/comment/commentstore"
	"salon_be/model/enrollment/enrollmentstore"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateCommentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input commentmodel.CreateComment

		if err := ctx.ShouldBind(&input); err != nil {
			panic(err)
		}

		requester := ctx.MustGet(common.CurrentUser).(common.Requester)

		db := appCtx.GetMainDBConnection()

		if err := db.Transaction(func(tx *gorm.DB) error {
			commentStore := commentstore.NewSQLStore(tx)
			enrollmentStore := enrollmentstore.NewSQLStore(tx)
			repo := commentrepo.NewCreateCommentRepo(commentStore, enrollmentStore)
			commentBusiness := commentbiz.NewCreateCommentBiz(repo)

			input.UserID = requester.GetUserId()
			if err := commentBusiness.CreateNewComment(ctx.Request.Context(), &input); err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
			return nil
		}); err != nil {
			panic(err)
		}
	}
}
