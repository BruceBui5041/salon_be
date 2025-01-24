package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/model/commission/commissionerror"
	"time"

	"gorm.io/gorm"
)

const CommissionEntityName = "Commission"

func init() {
	modelhelper.RegisterModel(Commission{})
}

type Commission struct {
	common.SQLModel `json:",inline"`
	Code            string     `json:"code" gorm:"column:code;uniqueIndex;not null;size:50"`
	PublishedAt     *time.Time `json:"published_at" gorm:"column:published_at;type:date"`
	RoleID          uint32     `json:"role_id" gorm:"column:role_id;not null"`
	Percentage      float64    `json:"percentage" gorm:"column:percentage;not null"`
	MinAmount       *float64   `json:"min_amount" gorm:"column:min_amount;"`
	MaxAmount       *float64   `json:"max_amount" gorm:"column:max_amount;"`
	Role            *Role      `json:"role,omitempty" gorm:"foreignKey:RoleID"`
}

func (Commission) TableName() string {
	return "commission"
}

func (c *Commission) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeCommission)
}

func (c *Commission) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	return
}

func (c *Commission) BeforeUpdate(tx *gorm.DB) (err error) {
	var existingCommission Commission
	if err := tx.Model(&Commission{}).
		Where("id = ?", c.Id).
		First(&existingCommission).
		Error; err != nil {
		return err
	}

	if existingCommission.PublishedAt != nil {
		return commissionerror.ErrCommissionPublished()
	}

	return nil
}
