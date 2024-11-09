package userbiz

import (
	"context"
	"errors"
	"regexp"
	"salon_be/common"
	"salon_be/component/hasher"
	"salon_be/component/tokenprovider"
	models "salon_be/model"
	"salon_be/model/auth/authconst"
	"salon_be/model/otp/otpmodel"
	"salon_be/model/user/usererror"
	"salon_be/model/user/usermodel"

	"github.com/jinzhu/copier"
)

type RegisterStorage interface {
	CreateNewUser(ctx context.Context, input *usermodel.CreateUser) error
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type OTPBiz interface {
	CreateOTP(ctx context.Context, data *otpmodel.CreateOTPInput) error
}

type registerBiz struct {
	registerStorage RegisterStorage
	hasher          hasher.Hasher
	tokenProvider   tokenprovider.Provider
	otpBiz          OTPBiz
}

func NewRegisterBusiness(
	registerStorage RegisterStorage,
	hasher hasher.Hasher,
	tokenProvider tokenprovider.Provider,
	otpBiz OTPBiz,
) *registerBiz {
	return &registerBiz{
		registerStorage: registerStorage,
		hasher:          hasher,
		tokenProvider:   tokenProvider,
		otpBiz:          otpBiz,
	}
}

// RegisterUser handles the registration of a new user
func (registerBiz *registerBiz) RegisterUser(
	ctx context.Context,
	inputData *usermodel.CreateUser,
	tokenExpiry int,
) (*usermodel.RegisterResponse, *models.User, error) {
	// Validate required fields
	if inputData.FirstName == "" {
		return nil, nil, usererror.ErrUserMissionRequireField(errors.New("firstname is required"))
	}

	if inputData.LastName == "" {
		return nil, nil, usererror.ErrUserMissionRequireField(errors.New("lastname is required"))
	}

	if inputData.AuthType == "" {
		return nil, nil, usererror.ErrUserMissionRequireField(errors.New("auth_type is required"))
	}

	if inputData.AuthType == authconst.AuthTypeEmail {
		if inputData.Email == "" {
			return nil, nil, usererror.ErrUserMissionRequireField(errors.New("email is required"))
		}
		// Validate email format
		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(inputData.Email) {
			return nil, nil, common.ErrInvalidRequest(errors.New("invalid email format"))
		}
	} else if inputData.AuthType == authconst.AuthTypePhone {
		if inputData.PhoneNumber == "" {
			return nil, nil, usererror.ErrUserMissionRequireField(errors.New("phonenumber is required"))
		}
	}

	switch inputData.AuthType {
	case authconst.AuthTypePassword:
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
	case authconst.AuthTypePhone:
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

	if inputData.AuthType == authconst.AuthTypePhone {
		err := registerBiz.otpBiz.CreateOTP(ctx, &otpmodel.CreateOTPInput{UserID: user.Id})
		if err != nil {
			return nil, nil, common.ErrInternal(err)
		}
	}

	payload := tokenprovider.TokenPayload{
		UserId:    int(user.Id),
		Roles:     user.Roles,
		Challenge: "otp",
	}

	accessToken, err := registerBiz.tokenProvider.Generate(payload, tokenExpiry)
	if err != nil {
		return nil, nil, common.ErrInternal(err)
	}

	var userRes usermodel.GetUserResponse
	copier.Copy(&userRes, user)
	return &usermodel.RegisterResponse{
		Token:     accessToken,
		User:      userRes,
		Challenge: "otp",
	}, user, nil
}
