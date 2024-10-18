package apihandler

import (
	"net/http"
	"video_server/common"
	"video_server/component"

	"github.com/gin-gonic/gin"
)

func DecodeUID(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		id := uid.GetLocalID()

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(id))
	}
}
