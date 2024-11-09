package userbiz

import (
	"context"
	"salon_be/common"
	"salon_be/component/hasher"
	"salon_be/component/logger"
	"salon_be/component/tokenprovider"
	models "salon_be/model"
	"salon_be/model/user/usererror"
	"salon_be/model/user/usermodel"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type LoginStorage interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

// type TokenConfig interface {
// 	GetAtExp() int
// 	GetRtExp() int
// }

type loginBusiness struct {
	loginStorage   LoginStorage
	tokenProvider  tokenprovider.Provider
	hasher         hasher.Hasher
	expiry         int
	localPublisher *gochannel.GoChannel
}

func NewLoginBusiness(
	storeUser LoginStorage,
	tokenProvicer tokenprovider.Provider,
	hasher hasher.Hasher,
	expiry int,
	localPublisher *gochannel.GoChannel,
) *loginBusiness {
	return &loginBusiness{
		loginStorage:   storeUser,
		tokenProvider:  tokenProvicer,
		hasher:         hasher,
		expiry:         expiry,
		localPublisher: localPublisher,
	}
}

func (business *loginBusiness) Login(ctx context.Context, data *usermodel.UserLogin) (*usermodel.LoginRes, error) {
	switch data.AuthType {
	case "phone_number":
		if data.PhoneNumber == "" {
			return nil, usererror.ErrUserPhoneNumberNotFound
		}
		return business.loginWithPhoneNumber(ctx, data)
	default:
		return business.loginWithEmailPassword(ctx, data)
	}
}

func (business *loginBusiness) loginWithPhoneNumber(ctx context.Context, data *usermodel.UserLogin) (*usermodel.LoginRes, error) {
	user, err := business.loginStorage.FindOne(
		ctx,
		map[string]interface{}{"phone_number": data.PhoneNumber},
		"Roles",
		"UserProfile",
	)

	if err != nil {
		logger.AppLogger.Error(ctx, "user not found", zap.Error(err))
		return nil, usererror.ErrUserPhoneNumberNotFound
	}

	payload := tokenprovider.TokenPayload{
		UserId:    int(user.Id),
		Roles:     user.Roles,
		Challenge: "otp",
	}

	accessToken, err := business.tokenProvider.Generate(payload, business.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	user.Mask(false)
	if err := watermill.PublishUserUpdated(
		ctx,
		business.localPublisher,
		&messagemodel.UserUpdatedMessage{UserId: user.GetFakeId()},
	); err != nil {
		return nil, common.ErrInternal(err)
	}

	var userRes usermodel.GetUserResponse
	copier.Copy(&userRes, user)
	return &usermodel.LoginRes{
		Token:     accessToken,
		User:      userRes,
		Challenge: "otp",
	}, nil
}

func (business *loginBusiness) loginWithEmailPassword(ctx context.Context, data *usermodel.UserLogin) (*usermodel.LoginRes, error) {
	user, err := business.loginStorage.FindOne(
		ctx,
		map[string]interface{}{"email": data.Email},
		"Roles",
		"UserProfile",
	)

	if err != nil {
		logger.AppLogger.Error(ctx, "user not found", zap.Error(err))
		return nil, usererror.ErrUsernameOrPasswordInvalid
	}

	pwdHashed := business.hasher.Hash(data.Password + user.Salt)
	if user.Password != pwdHashed {
		return nil, usererror.ErrUsernameOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserId: int(user.Id),
		Roles:  user.Roles,
	}

	accessToken, err := business.tokenProvider.Generate(payload, business.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	user.Mask(false)

	if err := watermill.PublishUserUpdated(
		ctx,
		business.localPublisher,
		&messagemodel.UserUpdatedMessage{UserId: user.GetFakeId()},
	); err != nil {
		return nil, common.ErrInternal(err)
	}

	// refreshToken, err := business.tokenProvider.Generate(payload, business.tokenConfig.GetRtExp())
	// if err != nil {
	// 	return nil, common.ErrInternal(err)
	// }

	// account := usermodel.NewAccount(accessToken, refreshToken)

	var userRes usermodel.GetUserResponse
	copier.Copy(&userRes, user)
	return &usermodel.LoginRes{
		Token: accessToken,
		User:  userRes,
	}, nil
}
