package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/tokenprovider/jwt"
	models "salon_be/model"
	"salon_be/utils"

	"github.com/gin-gonic/gin"
)

func ErrWrongAuthHeader(err error) *common.AppError {
	return common.NewCustomError(err, "wrong authen header", "ErrWrongAuthHeader")
}

func RequiredAuth(appCtx component.AppContext) func(ctx *gin.Context) {
	jwtProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie(appconst.AccessTokenName)
		if err != nil {
			panic(ErrWrongAuthHeader(errors.New("access_token cookie not found")))
		}

		payload, err := jwtProvider.Validate(token)
		if err != nil {
			panic(err)
		}

		isChallengeAPI := ctx.GetBool("isChallengeAPI")
		if !isChallengeAPI && payload.Challenge != "" {
			panic(ErrWrongAuthHeader(errors.New("have to passed authentication challenges")))
		}

		// Try to get user from cache
		appCache := appCtx.GetAppCache()

		userTypeObj := common.SQLModel{Id: uint32(payload.UserId)}
		userTypeObj.GenUID(common.DbTypeUser)

		cachedUser, err := appCache.GetUserCache(ctx.Request.Context(), userTypeObj.GetFakeId())
		if err != nil || cachedUser == "" {
			// clear access cookie
			utils.ClearServerJWTTokenCookie(ctx)
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

		ctx.Set(common.CurrentUser, user)
		SetCurrentUserToReqCtx(ctx, user)

		ctx.Next()
	}
}

func SetCurrentUserToReqCtx(c *gin.Context, user *models.User) {
	ctx := context.WithValue(c.Request.Context(), common.CurrentUser, user)
	c.Request = c.Request.WithContext(ctx)
}
