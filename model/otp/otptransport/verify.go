package otptransport

import (
	"errors"
	"net/http"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component"
	"salon_be/component/hasher"
	"salon_be/component/tokenprovider/jwt"
	"salon_be/model/otp/otpbiz"
	"salon_be/model/otp/otpmodel"
	"salon_be/model/otp/otprepo"
	"salon_be/model/otp/otpstore"
	"salon_be/model/user/userstore"
	"salon_be/utils"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func VerifyOTP(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input otpmodel.VerifyOTPInput
		if err := c.ShouldBind(&input); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester, ok := c.MustGet(common.CurrentUser).(common.Requester)
		if !ok {
			panic(common.ErrInvalidRequest(errors.New("requester not found")))
		}
		input.UserID = requester.GetUserId()

		db := appCtx.GetMainDBConnection()
		result := otpmodel.VerifyOTPResponse{}
		if err := db.Transaction(func(tx *gorm.DB) error {
			optStore := otpstore.NewSQLStore(tx)
			userStore := userstore.NewSQLStore(tx)
			tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
			md5 := hasher.NewMD5Hash()

			repo := otprepo.NewVerifyRepo(optStore, userStore)

			biz := otpbiz.NewVerifyBiz(
				repo,
				tokenProvider,
				md5,
				appconst.TokenExpiry,
			)

			res, err := biz.VerifyOTP(c.Request.Context(), &input)
			if err != nil {
				return err
			}

			err = copier.Copy(&result, res)
			if err != nil {
				return common.ErrInternal(err)
			}
			utils.WriteServerJWTTokenCookie(c, res.Token.Token)
			return nil
		}); err != nil {
			panic(err)
		}

		requester.Mask(false)
		if err := watermill.PublishUserUpdated(
			c.Request.Context(),
			appCtx.GetLocalPubSub().GetBlockPubSub(),
			&messagemodel.UserUpdatedMessage{UserId: requester.GetFakeId()},
		); err != nil {
			panic(common.ErrInternal(err))
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))

	}
}
