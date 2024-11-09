package usermodel

import (
	"salon_be/common"
	"salon_be/component/tokenprovider"
)

// CreateUser represents the data needed to create a new user
type CreateUser struct {
	*common.SQLModel  `json:",inline"`
	LastName          string `json:"lastname"`
	FirstName         string `json:"firstname"`
	Email             string `json:"email"`
	PhoneNumber       string `json:"phone_number"`
	Password          string `json:"password,omitempty"`
	Salt              string `json:"salt,omitempty"`
	AuthType          string `json:"auth_type"`
	AuthProviderID    string `json:"auth_provider_id,omitempty"`
	AuthProviderToken string `json:"auth_provider_token,omitempty"`
	ProfilePictureURL string `json:"profile_picture_url,omitempty"`
}

func (u *CreateUser) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

type RegisterResponse struct {
	Token     *tokenprovider.Token `json:"token"`
	User      GetUserResponse      `json:"user"`
	Challenge string               `json:"challenge"`
}
