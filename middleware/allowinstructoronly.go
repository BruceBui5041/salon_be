package middleware

import (
	"encoding/json"
	"errors"
	"salon_be/common"
	"salon_be/component"
	models "salon_be/model"

	"github.com/gin-gonic/gin"
)

func ErrNotAnIntructor() *common.AppError {
	return common.NewCustomError(errors.New("allow only instructor"), "allow only instructor", "ErrNotAnIntructor")
}

func AllowIntructorOnly(appCtx component.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("cannot find requester")))
		}
		var roles []models.Role
		err := json.Unmarshal(requester.GetRoles(c), &roles)
		if err != nil {
			panic(err)
		}

		for _, role := range roles {
			if role.Name == "instructor" {
				c.Next()
				return
			}
		}

		panic(ErrNotAnIntructor())
	}
}
