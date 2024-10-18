package commenttransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/comment/commentbiz"
	"video_server/model/comment/commentmodel"
	"video_server/model/comment/commentrepo"
	"video_server/model/comment/commentstore"
	"video_server/model/enrollment/enrollmentstore"

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
