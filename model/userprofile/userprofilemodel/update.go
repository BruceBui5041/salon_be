package userprofilemodel

import "mime/multipart"

type UpdateProfileModel struct {
	FirstName         *string               `json:"firstname,omitempty"`
	LastName          *string               `json:"lastname,omitempty"`
	PhoneNumber       *string               `json:"phone_number,omitempty"`
	Occupation        *string               `json:"occupation,omitempty"`
	Biography         *string               `json:"biography,omitempty"`
	LinkedIn          *string               `json:"linkedin,omitempty"`
	Facebook          *string               `json:"facebook,omitempty"`
	Twitter           *string               `json:"twitter,omitempty"`
	Instagram         *string               `json:"instagram,omitempty"`
	ProfilePictureURL *multipart.FileHeader `json:"profile_picture_url" form:"profile_picture_url"`
}
