package usermodel

import (
	"salon_be/common"
	models "salon_be/model"
	"time"
)

type GetUserResponse struct {
	common.SQLModel   `json:",inline"`
	Status            string                      `json:"status"`
	LastName          string                      `json:"lastname" gorm:"column:lastname;"`
	FirstName         string                      `json:"firstname" gorm:"column:firstname;"`
	Email             string                      `json:"email"`
	ProfilePictureURL string                      `json:"profile_picture_url"`
	Roles             []GetUserRoleResponse       `json:"roles"`
	Enrollments       []GetUserEnrollmentResponse `json:"-"`
	Auths             []models.UserAuth           `json:"-"`
}

func (u *GetUserResponse) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

type GetUserRoleResponse struct {
	common.SQLModel `json:",inline"`
	Name            string `json:"name"`
	Description     string
}

type GetUserEnrollmentResponse struct {
	common.SQLModel `json:",inline"`
	EnrolledAt      time.Time              `json:"enrolled_at"`
	Service         GetUserServiceResponse `json:"service"`
}

type GetUserServiceResponse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
	Description     string `json:"description"`
}
