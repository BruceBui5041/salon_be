package models

type UserRole struct {
	UserID uint `json:"user_id" gorm:"primaryKey;column:user_id;joinReferences:UserID"`
	RoleID uint `json:"role_id" gorm:"primaryKey;column:role_id;joinReferences:RoleID"`
}

func (UserRole) TableName() string {
	return "user_role"
}
