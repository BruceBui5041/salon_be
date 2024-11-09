package usermodel

import (
	"salon_be/component/tokenprovider"
)

type LoginRes struct {
	Token     *tokenprovider.Token `json:"token"`
	User      GetUserResponse      `json:"user"`
	Challenge string               `json:"challenge"`
}
