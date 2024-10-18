package middleware

import (
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recover(ac component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					c.AbortWithStatusJSON(appErr.StatusCode, appErr)
					// NOTE: panic lại để gin có thể log ra lại được stack trace của err
					logger.AppLogger.Error(c.Request.Context(), "gin recover by", zap.Error(err.(error)))
					panic(err)
					// return
				}

				appErr := common.ErrInternal(err.(error))
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				logger.AppLogger.Error(c.Request.Context(), "gin recover by", zap.Error(err.(error)))
				panic(err)
				// return
			}
		}()

		c.Next()
	}
}
