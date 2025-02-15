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
	LastName        string        `json:"lastname" gorm:"column:lastname;"`
	FirstName       string        `json:"firstname" gorm:"column:firstname;"`
	Email           string        `json:"email" gorm:"column:email;uniqueIndex;size:100"`
	PhoneNumber     string        `json:"phone_number" gorm:"column:phone_number;uniqueIndex;size:20"`
	Roles           []*Role       `json:"roles" gorm:"many2many:user_role;joinForeignKey:UserID;joinReferences:RoleID"`
	Auths           []UserAuth    `json:"auths" gorm:"foreignKey:UserID"`
	CreatedServices []Service     `json:"created_services" gorm:"foreignKey:CreatorID"`
	Enrollments     []*Enrollment `json:"enrollments" `
	Progress        []Progress    `json:"progress"`
	Salt            string        `json:"-" gorm:"column:salt;"`
	Password        string        `json:"-" gorm:"column:password;"`
	UserProfile     *UserProfile  `json:"user_profile,omitempty" gorm:"foreignKey:UserID"`
	Comments        []*Comment    `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	OTPs            []*OTP        `json:"otp,omitempty" gorm:"foreignKey:UserID"`
	UserDevice      *UserDevice   `json:"user_device,omitempty" gorm:"foreignKey:UserID"`

	GroupProviders []*GroupProvider `json:"group_providers,omitempty" gorm:"many2many:m2m_group_provider_users;foreignKey:Id;joinForeignKey:UserID;References:Id;joinReferences:GroupProviderID"`
	AdminGroups    []*GroupProvider `json:"admin_groups,omitempty" gorm:"foreignKey:AdminID"`
	CreatedGroups  []*GroupProvider `json:"created_groups,omitempty" gorm:"foreignKey:CreatorID"`

	OwnedCertificates   []*Certificate `json:"owned_certificates" gorm:"foreignKey:OwnerID"`
	CreatedCertificates []*Certificate `json:"created_certificates" gorm:"foreignKey:CreatorID"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Status = "inactive"
	if u.Email == "" && u.PhoneNumber != "" {
		u.Email = u.PhoneNumber
	}

	return nil
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
		return role.Code == "ADMIN" || role.Code == "SUPER_ADMIN"
	})

	return has
}

func (u *User) IsUser() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Code == "USER"
	})

	return has
}

func (u *User) IsProvider() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Code == "PROVIDER"
	})

	return has
}

func (u *User) IsGroupProviderAdmin() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Code == "GROUP_PROVIDER_ADMIN"
	})

	return has
}

func (u *User) IsSuperAdmin() bool {
	_, has := lo.Find(u.Roles, func(role *Role) bool {
		return role.Code == "SUPER_ADMIN"
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
