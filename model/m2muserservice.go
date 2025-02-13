package models

type UserService struct {
	UserID           uint32 `json:"user_id" gorm:"primaryKey;column:user_id;joinReferences:UserID"`
	ServiceVersionID uint32 `json:"service_version_id" gorm:"primaryKey;column:service_version_id;joinReferences:ServiceVersionID"`
}

func (UserService) TableName() string {
	return "user_service"
}
