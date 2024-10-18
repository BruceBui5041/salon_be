package coursebiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type GetCourseByIDStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type getCourseByIDBiz struct {
	store GetCourseByIDStore
}

func NewGetCourseByIDBiz(store GetCourseByIDStore) *getCourseByIDBiz {
	return &getCourseByIDBiz{store: store}
}

func (biz *getCourseByIDBiz) GetCourseByID(
	ctx context.Context,
	id int,
	moreInfo ...string,
) (*coursemodel.CourseResponse, error) {
	course, err := biz.store.FindOne(
		ctx,
		map[string]interface{}{"id": id},
		"Category",
		"Creator.UserProfile",
		"IntroVideo",
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
		common.PreloadInfo{
			Name: "Comments",
			Function: func(d *gorm.DB) *gorm.DB {
				return d.Limit(5).Order("`rate` DESC")
			},
		},
		"Comments.User.UserProfile",
	)
	if err != nil {
		if err == common.RecordNotFound {
			return nil, common.ErrCannotGetEntity(models.CourseEntityName, err)
		}
		return nil, common.ErrCannotGetEntity(models.CourseEntityName, err)
	}

	var getCourse coursemodel.CourseResponse
	copier.Copy(&getCourse, course)

	return &getCourse, nil
}
