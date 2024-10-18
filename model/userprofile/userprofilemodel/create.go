package userprofilemodel

import (
	"salon_be/common"
)

type CreateUserProfile struct {
	common.SQLModel `json:",inline"`
	UserID          uint32 `json:"user_id"`
	FirstName       string `json:"firstname"`
	LastName        string `json:"lastname"`
	PhoneNumber     string `json:"phone_number"`
	Occupation      string `json:"occupation"`
	Biography       string `json:"biography"`
	LinkedIn        string `json:"linkedin"`
	Facebook        string `json:"facebook"`
	Twitter         string `json:"twitter"`
	Instagram       string `json:"instagram"`
}

func (cup *CreateUserProfile) Mask(isAdmin bool) {
	cup.GenUID(common.DBTypeUserProfile)
}
