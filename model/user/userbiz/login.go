package userbiz

import (
	"context"
	"video_server/common"
	"video_server/component/hasher"
	"video_server/component/tokenprovider"
	models "video_server/model"
	"video_server/model/user/usermodel"
	"video_server/watermill"
	"video_server/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/jinzhu/copier"
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
	user, err := business.loginStorage.FindOne(
		ctx,
		map[string]interface{}{"email": data.Email},
		"Roles",
		"Enrollments.Course",
		"UserProfile",
	)

	if err != nil {
		return nil, usermodel.ErrUsernameOrPasswordInvalid
	}

	pwdHashed := business.hasher.Hash(data.Password + user.Salt)
	if user.Password != pwdHashed {
		return nil, usermodel.ErrUsernameOrPasswordInvalid
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
