package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"video_server/common"
	"video_server/component/genericapi/modelhelper"
	"video_server/storagehandler"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

const CourseEntityName = "Course"

func init() {
	modelhelper.RegisterModel(Course{})
}

type ReviewInfo struct {
	Stars uint8 `json:"stars"`
	Count uint  `json:"count"`
}

type Course struct {
	common.SQLModel `json:",inline"`
	Title           string              `json:"title" gorm:"column:title;not null;size:255"`
	Description     string              `json:"description" gorm:"column:description;type:text"`
	CreatorID       uint32              `json:"creator_id" gorm:"column:creator_id;index"`
	CategoryID      uint32              `json:"category_id" gorm:"column:category_id;index"`
	IntroVideoID    *uint32             `json:"intro_video_id,omitempty" gorm:"column:intro_video_id;index"`
	Creator         *User               `json:"creator,omitempty" gorm:"constraint:OnDelete:SET NULL;"`
	Category        *Category           `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL;"`
	Instructors     []User              `json:"instructors" gorm:"many2many:user_course;"`
	Lessons         []Lesson            `json:"lessons" gorm:"foreignKey:CourseID"`
	Lectures        []Lecture           `json:"lectures" gorm:"foreignKey:CourseID"`
	Videos          []Video             `json:"videos" gorm:"foreignKey:CourseID"`
	IntroVideo      *Video              `json:"intro_video,omitempty" gorm:"foreignKey:IntroVideoID"`
	Enrollments     []Enrollment        `json:"enrollments,omitempty" gorm:"foreignKey:CourseID"`
	Slug            string              `json:"slug" gorm:"column:slug;not null;size:255"`
	Thumbnail       string              `json:"thumbnail" gorm:"column:thumbnail;type:text"`
	Overview        string              `json:"overview" gorm:"column:overview;type:text"`
	Price           decimal.Decimal     `json:"price" gorm:"column:price;type:decimal(10,2);not null"`
	DiscountedPrice decimal.NullDecimal `json:"discounted_price" gorm:"column:discounted_price;type:decimal(10,2);"`
	DifficultyLevel string              `json:"difficulty_level" gorm:"column:difficulty_level;type:ENUM('beginner','intermediate','advanced','expert');not null"`
	Comments        []*Comment          `json:"comments,omitempty" gorm:"foreignKey:CourseID"`
	LessonCount     uint16              `json:"lesson_count" gorm:"column:lesson_count;type:unsigned int"`
	StudentCount    uint16              `json:"student_count" gorm:"column:student_count;type:unsigned int"`
	RatingCount     uint                `json:"rating_count" gorm:"column:rating_count;type:unsigned int"`
	ReviewInfo      ReviewInfos         `json:"review_info" gorm:"column:review_info;type:json"`
	AverageRating   decimal.Decimal     `json:"avg_rating" gorm:"column:avg_rating;type:decimal(3,1)"`
}

func (Course) TableName() string {
	return "course"
}

func (c *Course) Mask(isAdmin bool) {
	c.GenUID(common.DbTypeCourse)
}

func (c *Course) AfterFind(tx *gorm.DB) (err error) {
	c.Mask(false)
	c.Thumbnail = storagehandler.AddPublicCloudFrontDomain(c.Thumbnail)
	return
}

type ReviewInfos []ReviewInfo

func (r *ReviewInfos) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, r)
}

func (r ReviewInfos) Value() (driver.Value, error) {
	return json.Marshal(r)
}
