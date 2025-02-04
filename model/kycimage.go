package models

import (
	"salon_be/common"
)

const KYCImageEntityName = "ImageUpload"

type KYCImageUpload struct {
	common.SQLModel `json:",inline"`
	UserID          uint32  `json:"user_id" gorm:"column:user_id;"`
	User            *User   `json:"-" gorm:"foreignKey:UserID;"`
	FileName        string  `json:"fileName" gorm:"column:file_name;"`
	Title           string  `json:"title" gorm:"column:title;"`
	Description     string  `json:"description" gorm:"column:description;"`
	Hash            string  `json:"hash" gorm:"column:hash;"`
	FileType        string  `json:"fileType" gorm:"column:file_type;"`
	UploadedDate    string  `json:"uploadedDate" gorm:"column:uploaded_date;"`
	StorageType     string  `json:"storageType" gorm:"column:storage_type;"`
	TokenId         string  `json:"tokenId" gorm:"column:token_id;"`
	Provider        string  `json:"provider" gorm:"column:provider;"`
	ImageID         *uint32 `json:"-" gorm:"column:image_id"`
	Image           *Image  `json:"image,omitempty" gorm:"foreignKey:ImageID"`
}

func (KYCImageUpload) TableName() string {
	return "kyc_image"
}
