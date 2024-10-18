package models

import (
	"context"
	"encoding/json"
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/component/logger"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const UserEntityName = "User"

func init() {
	modelhelper.RegisterModel(User{})
}

type User struct {
	common.SQLModel `json:",inline"`
	LastName        string           `json:"lastname" gorm:"column:lastname;"`
	FirstName       string           `json:"firstname" gorm:"column:firstname;"`
	Email           string           `json:"email" gorm:"column:email;uniqueIndex;not null;size:100"`
	Roles           []*Role          `json:"roles" gorm:"many2many:user_role;joinForeignKey:UserID;joinReferences:RoleID"`
	Auths           []UserAuth       `json:"auths" gorm:"foreignKey:UserID"`
	CreatedServices []ServiceVersion `json:"created_services" gorm:"foreignKey:CreatorID"`
	Enrollments     []*Enrollment    `json:"enrollments" `
	Progress        []Progress       `json:"progress"`
	Salt            string           `json:"-" gorm:"column:salt;"`
	Password        string           `json:"-" gorm:"column:password;"`
	UserProfile     *UserProfile     `json:"user_profile,omitempty" gorm:"foreignKey:UserID"`
	Comments        []*Comment       `json:"comments,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

func (u *User) GetUserId() uint32 {
	return u.Id
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) IsAdmin() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Name == "admin"
	})

	return has
}

func (u *User) IsStudent() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Name == "student"
	})

	return has
}

func (u *User) IsInstructor() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Name == "instructor"
	})

	return has
}

func (u *User) IsSuperAdmin() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Name == "super_admin"
	})

	return has
}

func (u *User) GetRoles(ctx context.Context) []byte {
	data, err := json.Marshal(u.Roles)
	if err != nil {
		logger.AppLogger.Error(ctx, "cannot marshal users roles", zap.Error(err))
		return nil
	}
	return data
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.Mask(false)
	return
}
