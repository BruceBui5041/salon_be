package usermodel

import (
	"salon_be/component/tokenprovider"
	models "salon_be/model"
)

type UserLogin struct {
	AuthType    string `json:"auth_type" form:"auth_type"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Email       string `json:"email" form:"email"`
	Password    string `json:"password" form:"password"`
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
