package models

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"

	"gorm.io/gorm"
)

const CommentEntityName = "Comment"

func init() {
	modelhelper.RegisterModel(Comment{})
}

type Comment struct {
	common.SQLModel `json:",inline"`
	UserID          uint32  `json:"user_id" gorm:"column:user_id;index"`
	CourseID        uint32  `json:"course_id" gorm:"column:course_id;index"`
	Content         string  `json:"content" gorm:"column:content;type:text"`
	Rate            uint8   `json:"rate" gorm:"column:rate;type:decimal(2,1);not null"`
	User            *User   `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`
	Course          *Course `json:"course,omitempty" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL;"`
}

func (Comment) TableName() string {
	return "comment"
}

func (c *Comment) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeComment)
}

func (c *Comment) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	return
}
