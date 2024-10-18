package apihandler

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"

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
