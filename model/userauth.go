package models

type UserAuth struct {
	ID                uint32 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID            uint32 `gorm:"column:user_id;index"`
	AuthType          string `gorm:"column:auth_type;not null;size:20"`
	AuthProviderID    string `gorm:"column:auth_provider_id;size:255"`
	AuthProviderToken string `gorm:"column:auth_provider_token"`
	User              User   `gorm:"constraint:OnDelete:CASCADE;"`
}

func (UserAuth) TableName() string {
	return "user_auth"
}
