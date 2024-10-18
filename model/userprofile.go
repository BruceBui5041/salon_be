package models

import (
	"salon_be/common"
	"salon_be/storagehandler"

	"gorm.io/gorm"
)

const UserProfileEntityName = "UserProfile"

type UserProfile struct {
	common.SQLModel   `json:",inline"`
	UserID            uint32 `json:"-" gorm:"column:user_id;uniqueIndex"`
	FirstName         string `json:"firstname" gorm:"column:firstname;size:50"`
	LastName          string `json:"lastname" gorm:"column:lastname;size:50"`
	User              *User  `json:"-" gorm:"foreignKey:UserID"`
	ProfilePictureURL string `json:"profile_picture_url" gorm:"column:profile_picture_url;size:255"`
	PhoneNumber       string `json:"phone_number" gorm:"column:phone_number;size:20"`
	Occupation        string `json:"occupation" gorm:"column:occupation;size:100"`
	Biography         string `json:"biography" gorm:"column:biography;type:text"`
	LinkedIn          string `json:"linkedin" gorm:"column:linkedin;size:255"`
	Facebook          string `json:"facebook" gorm:"column:facebook;size:255"`
	Twitter           string `json:"twitter" gorm:"column:twitter;size:255"`
	Instagram         string `json:"instagram" gorm:"column:instagram;size:255"`
}

func (UserProfile) TableName() string {
	return "user_profile"
}

func (up *UserProfile) Mask(isAdmin bool) {
	up.GenUID(common.DBTypeUserProfile)
}

func (up *UserProfile) AfterFind(tx *gorm.DB) (err error) {
	up.Mask(false)
	up.ProfilePictureURL = storagehandler.AddPublicCloudFrontDomain(up.ProfilePictureURL)
	return
}
