package ws

import (
	"encoding/json"
	"errors"
	"net/http"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/logger"
	"salon_be/component/tokenprovider/jwt"
	models "salon_be/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func wsAuthorize(c *gin.Context, appCtx component.AppContext) *models.User {
	ctx := c.Request.Context()
	accessToken, err := c.Cookie(appconst.AccessTokenName)
	if err != nil {
		logger.AppLogger.Error(ctx, "Missing access token cookie",
			zap.Error(err),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
		logger.AppLogger.Error(ctx, "missing access token cookie")
		return nil
	}
	jwtProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

	payload, err := jwtProvider.Validate(accessToken)
	if err != nil {
		panic(err)
	}

	if payload.Challenge != "" {
		panic(errors.New("have to passed authentication challenges"))
	}

	// Try to get user from cache
	appCache := appCtx.GetAppCache()

	userTypeObj := common.SQLModel{Id: uint32(payload.UserId)}
	userTypeObj.GenUID(common.DbTypeUser)

	cachedUser, err := appCache.GetUserCache(ctx, userTypeObj.GetFakeId())
	if err != nil || cachedUser == "" {
		panic(common.ErrNoPermission(errors.New("token expired")))
	}

	var user *models.User
	// User found in cache, unmarshal it
	err = json.Unmarshal([]byte(cachedUser), &user)
	if err != nil {
		// If there's an error unmarshalling, we'll fetch from the database
		user = nil
	}

	if payload.Challenge == "" && user.Status != common.StatusActive {
		panic(common.ErrNoPermission(errors.New("account unavailable")))
	}

	uid, err := common.FromBase58(user.GetFakeId())
	if err != nil {
		panic(err)
	}

	user.Id = uid.GetLocalID()

	return user
}
