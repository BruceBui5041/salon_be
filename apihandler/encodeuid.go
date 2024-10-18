package apihandler

import (
	"net/http"
	"salon_be/common"
	"salon_be/component"
	"strconv"

	"github.com/gin-gonic/gin"
)

func EncodeUID(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			common.ErrInvalidRequest(err)
		}

		dbType, err := strconv.Atoi(c.Param("dbtype"))
		if err != nil {
			common.ErrInvalidRequest(err)
		}

		obj := common.SQLModel{Id: uint32(id)}
		obj.GenUID(dbType)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(obj.GetFakeId()))
	}
}
