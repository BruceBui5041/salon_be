package models

type UserRole struct {
	UserID uint32 `json:"user_id" gorm:"primaryKey;column:user_id;joinReferences:UserID"`
	RoleID uint32 `json:"role_id" gorm:"primaryKey;column:role_id;joinReferences:RoleID"`
}

func (UserRole) TableName() string {
	return "user_role"
}
