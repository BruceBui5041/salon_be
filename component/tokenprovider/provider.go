package tokenprovider

import (
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"time"
)

type Provider interface {
	Generate(data TokenPayload, expiry int) (*Token, error)
	Validate(token string) (*TokenPayload, error)
}

type Token struct {
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
	Expiry  int       `json:"expiry"`
}

type TokenPayload struct {
	UserId    int            `json:"user_id"`
	Roles     []*models.Role `json:"roles"`
	Challenge string         `json:"challenge"`
}

var (
	ErrNotFound      = common.NewCustomError(errors.New("token not found"), "token not found", "ErrNotFound")
	ErrEncodingToken = common.NewCustomError(errors.New("error encoding token"), "error encoding token", "ErrEncodingToken")
	ErrInvalidToken  = common.NewCustomError(errors.New("invalid token"), "invalid token", "ErrInvalidTokne")
)
