package categorymodel

import (
	"salon_be/common"
	"salon_be/component/genericapi/modelhelper"
	models "salon_be/model"

	"github.com/samber/lo"
)

func init() {
	modelhelper.RegisterResponseType(models.Category{}.TableName(), CategoryResponse{})
}

type CategoryResponse struct {
	common.SQLModel `json:",inline"`
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	Courses         []models.Course `json:"course,omitempty"`
	CourseCount     int             `json:"course_count"`
}

func (cr *CategoryResponse) CountCourse() {
	cr.CourseCount = lo.CountBy(cr.Courses, func(course models.Course) bool {
		return course.Status == "active"
	})
}

func (cr *CategoryResponse) RemoveCoursesResponse() {
	cr.Courses = []models.Course{}
}

// func (cr *CategoryResponse) Mask(isAdmin bool) {
// 	cr.GenUID(common.DBTypeCategory)
// }
