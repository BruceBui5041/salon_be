package usermodel

import (
	"errors"
	"salon_be/common"
	"salon_be/component/tokenprovider"
	models "salon_be/model"
)

type UserLogin struct {
	Email    string `json:"email" form:"email" gorm:"column:email;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return models.User{}.TableName()
}

type Account struct {
	AccessToken  *tokenprovider.Token `json:"access_token"`
	RefreshToken *tokenprovider.Token `json:"refresh_token"`
}

func NewAccount(atok, rtok *tokenprovider.Token) *Account {
	return &Account{
		AccessToken:  atok,
		RefreshToken: rtok,
	}
}

var (
	ErrUsernameOrPasswordInvalid = common.NewCustomError(
		errors.New("username or password is invalid"),
		"username or password is invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrEmailIsAlreadyExisted = common.NewCustomError(
		errors.New("email is alread existed"),
		"email is alread existed",
		"ErrEmailIsAlreadyExisted",
	)
)
