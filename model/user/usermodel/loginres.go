package usermodel

import (
	"video_server/component/tokenprovider"
)

type LoginRes struct {
	Token *tokenprovider.Token `json:"token"`
	User  GetUserResponse      `json:"user"`
}
