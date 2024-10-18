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
	Enrollments       []GetUserEnrollmentResponse `json:"enrollments"`
	Auths             []models.UserAuth           `json:"auths"`
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
	EnrolledAt      time.Time             `json:"enrolled_at"`
	Course          GetUserCourseResponse `json:"course"`
}

type GetUserCourseResponse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title"`
	Description     string `json:"description"`
}
