package models

type M2MServiceVersionImage struct {
	ServiceVersionID uint32          `gorm:"column:service_version_id;primaryKey;not null;index;type:uint"`
	ImageID          uint32          `gorm:"column:image_id;primaryKey;not null;index;type:uint"`
	ServiceVersion   *ServiceVersion `gorm:"foreignKey:ServiceVersionID;references:Id;constraint:OnDelete:CASCADE;"`
	Image            *Image          `gorm:"foreignKey:ImageID;references:Id;constraint:OnDelete:CASCADE;"`
}

func (M2MServiceVersionImage) TableName() string {
	return "service_version_images"
}
