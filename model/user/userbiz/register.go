package userbiz

import (
	"context"
	"errors"
	"regexp"
	"salon_be/common"
	"salon_be/component/hasher"
	"salon_be/component/tokenprovider"
	models "salon_be/model"
	"salon_be/model/user/usermodel"
)

type RegisterStorage interface {
	CreateNewUser(ctx context.Context, input *usermodel.CreateUser) error
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type registerBiz struct {
	registerStorage RegisterStorage
	hasher          hasher.Hasher
	tokenProvider   tokenprovider.Provider
}

func NewRegisterBusiness(
	registerStorage RegisterStorage,
	hasher hasher.Hasher,
	tokenProvider tokenprovider.Provider,
) *registerBiz {
	return &registerBiz{
		registerStorage: registerStorage,
		hasher:          hasher,
		tokenProvider:   tokenProvider,
	}
}

// RegisterUser handles the registration of a new user
func (registerBiz *registerBiz) RegisterUser(
	ctx context.Context,
	inputData *usermodel.CreateUser,
	tokenExpiry int,
) (*tokenprovider.Token, *models.User, error) {
	// Validate required fields
	if inputData.FirstName == "" {
		return nil, nil, errors.New("first name is required")
	}
	if inputData.LastName == "" {
		return nil, nil, errors.New("last name is required")
	}

	if inputData.Email == "" {
		return nil, nil, errors.New("email is required")
	}

	if inputData.Password != "" {
		inputData.AuthType = "password"
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(inputData.Email) {
		return nil, nil, errors.New("invalid email format")
	}

	// Validate auth type and corresponding fields
	switch inputData.AuthType {
	case "password":
		if inputData.Password == "" {
			return nil, nil, errors.New("password is required for password auth type")
		}
		// Hash password
		salt := common.GenSalt(50)
		hashedPassword := registerBiz.hasher.Hash(inputData.Password + salt)
		inputData.Salt = salt
		inputData.Password = string(hashedPassword)
	case "oauth":
		if inputData.AuthProviderID == "" || inputData.AuthProviderToken == "" {
			return nil, nil, errors.New("auth provider ID and token are required for oauth auth type")
		}
	default:
		return nil, nil, errors.New("invalid auth type")
	}

	err := registerBiz.registerStorage.CreateNewUser(ctx, inputData)

	if err != nil {
		return nil, nil, err
	}

	user, err := registerBiz.registerStorage.FindOne(ctx, map[string]interface{}{"email": inputData.Email})
	if err != nil {
		return nil, nil, common.ErrInternal(err)
	}

	user.Mask(false)

	payload := tokenprovider.TokenPayload{
		UserId: int(user.Id),
		Roles:  user.Roles,
	}

	accessToken, err := registerBiz.tokenProvider.Generate(payload, tokenExpiry)
	if err != nil {
		return nil, nil, common.ErrInternal(err)
	}

	return accessToken, user, nil
}
