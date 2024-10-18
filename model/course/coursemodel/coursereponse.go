package coursemodel

import (
	"video_server/common"
	"video_server/component/genericapi/modelhelper"
	models "video_server/model"
	"video_server/utils/customtypes"

	"github.com/shopspring/decimal"
)

func init() {
	modelhelper.RegisterResponseType(models.Course{}.TableName(), CourseResponse{})
}

type CourseResponse struct {
	common.SQLModel `json:",inline"`
	Title           string                        `json:"title"`
	Description     string                        `json:"description"`
	Creator         *models.User                  `json:"creator,omitempty"`
	Category        *models.Category              `json:"category,omitempty"`
	IntroVideo      *CourseVideoReponse           `json:"intro_video,omitempty"`
	Enrollments     []models.Enrollment           `json:"enrollments,omitempty"`
	Lessons         []CourseLessonResponse        `json:"lessons"`
	Lectures        []CourseLectureResponse       `json:"lectures,omitempty"`
	Slug            string                        `json:"slug"`
	Thumbnail       string                        `json:"thumbnail"`
	Price           customtypes.DecimalString     `json:"price" `
	DiscountedPrice customtypes.NullDecimalString `json:"discounted_price" `
	DifficultyLevel string                        `json:"difficulty_level"`
	Overview        string                        `json:"overview"`
	Comments        []*CommentResponse            `json:"comments,omitempty"`
	RatingCount     uint                          `json:"rating_count"`
	LessonCount     uint16                        `json:"lesson_count"`
	StudentCount    uint16                        `json:"student_count"`
	ReviewInfo      models.ReviewInfos            `json:"review_info"`
	AverageRating   decimal.Decimal               `json:"avg_rating"`
}

func (g *CourseResponse) Mask(isAdmin bool) {
	g.GenUID(common.DbTypeCourse)
}

type CommentResponse struct {
	common.SQLModel `json:",inline"`
	Content         string        `json:"content"`
	Rate            uint8         `json:"rate"`
	User            *UserResponse `json:"user,omitempty"`
}

func (c *CommentResponse) Mask(isAdmin bool) {
	c.GenUID(common.DBTypeComment)
}

type UserResponse struct {
	common.SQLModel   `json:",inline"`
	Status            string `json:"status"`
	LastName          string `json:"lastname"`
	FirstName         string `json:"firstname"`
	Email             string `json:"email"`
	ProfilePictureURL string `json:"profile_picture_url"`
}

func (u *UserResponse) Mask(isAdmin bool) {
	u.GenUID(common.DbTypeUser)
}

type CourseVideoReponse struct {
	common.SQLModel `json:",inline"`
	Title           string `json:"title" `
	Description     string `json:"description"`
	// VideoURL        string `json:"video_url" `
	ThumbnailURL string `json:"thumbnail_url" `
	Duration     int    `json:"duration"`
	Order        int    `json:"order"`
	AllowPreview bool   `json:"allow_preview" `
	Overview     string `json:"overview"`
}

func (g *CourseVideoReponse) Mask(isAdmin bool) {
	g.GenUID(common.DbTypeVideo)
}

type CourseLessonResponse struct {
	common.SQLModel `json:",inline"`
	CourseID        uint32               `json:"-"`
	Title           string               `json:"title"`
	Description     string               `json:"description"`
	Duration        int                  `json:"duration"`
	Order           int                  `json:"order"`
	Videos          []CourseVideoReponse `json:"videos"`
	Type            string               `json:"type"`
	AllowPreview    bool                 `json:"allow_preview"`
}

func (CourseLessonResponse) TableName() string {
	return models.Lesson{}.TableName()
}

func (l *CourseLessonResponse) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLesson)
}

type CourseLectureResponse struct {
	common.SQLModel `json:",inline"`
	CourseID        uint32                 `json:"-" gorm:"index"`
	Title           string                 `json:"title" gorm:"not null;size:255"`
	Description     string                 `json:"description" gorm:"type:text"`
	Order           int                    `json:"order"`
	Course          models.Course          `json:"-" gorm:"foreignKey:CourseID;constraint:OnDelete:SET NULL;"`
	Lessons         []CourseLessonResponse `json:"lessons" gorm:"foreignKey:LectureID"`
	Duration        int                    `json:"duration"`
}

func (l *CourseLectureResponse) Mask(isAdmin bool) {
	l.GenUID(common.DBTypeLecture)
}
