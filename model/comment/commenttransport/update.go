package commenttransport

import (
	"net/http"
	"video_server/common"
	"video_server/component"
	"video_server/model/comment/commentbiz"
	"video_server/model/comment/commentmodel"
	"video_server/model/comment/commentrepo"
	"video_server/model/comment/commentstore"

	"github.com/gin-gonic/gin"
)

func UpdateCommentHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		var input commentmodel.UpdateComment
		if err := ctx.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()

		commentStore := commentstore.NewSQLStore(db)
		repo := commentrepo.NewUpdateCommentRepo(commentStore)
		commentBusiness := commentbiz.NewUpdateCommentBiz(repo)

		if err := commentBusiness.UpdateComment(ctx.Request.Context(), uint32(id), &input); err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse("ok"))
	}
}
