package coursebiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ListCourseStore interface {
	FindAll(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) ([]models.Course, error)
}

type listCourseBiz struct {
	listCourseStore ListCourseStore
}

func NewCourseBiz(listCourseStore ListCourseStore) *listCourseBiz {
	return &listCourseBiz{listCourseStore: listCourseStore}
}

func (biz *listCourseBiz) ListCourses(ctx context.Context) ([]coursemodel.CourseResponse, error) {
	courses, err := biz.listCourseStore.FindAll(
		ctx,
		map[string]interface{}{"status": "active"},
		"Creator.UserProfile",
		"Category",
		common.PreloadInfo{
			Name: "Lectures",
			Function: func(d *gorm.DB) *gorm.DB {
				return d.Order("`order` ASC").Order("`id` ASC")
			},
		},
		common.PreloadInfo{
			Name: "Lectures.Lessons",
			Function: func(d *gorm.DB) *gorm.DB {
				return d.Order("`order` ASC").Order("`id` ASC")
			},
		},
		"Lectures.Lessons.Videos",
	)

	if err != nil {
		return nil, common.ErrCannotListEntity(models.CourseEntityName, err)
	}

	var courseRes []coursemodel.CourseResponse
	copier.Copy(&courseRes, courses)

	return courseRes, nil
}
