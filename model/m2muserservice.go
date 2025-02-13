package models

func init() {
	// modelhelper.RegisterModel(UserService{})
}

type UserService struct {
	UserID           uint32          `json:"-" gorm:"primaryKey;column:user_id"`
	ServiceVersionID uint32          `json:"-" gorm:"primaryKey;column:service_version_id"`
	User             *User           `json:"user,omitempty" gorm:"foreignKey:UserID;references:Id"`
	ServiceVersion   *ServiceVersion `json:"service_version,omitempty" gorm:"foreignKey:ServiceVersionID;references:Id"`
}

func (UserService) TableName() string {
	return "user_service"
}
