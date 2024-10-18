package coursebiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"

	"github.com/jinzhu/copier"
)

type GetCourseVideosStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type getCourseVideosBiz struct {
	store GetCourseVideosStore
}

func NewGetCourseVideosBiz(store GetCourseVideosStore) *getCourseVideosBiz {
	return &getCourseVideosBiz{store: store}
}

func (biz *getCourseVideosBiz) GetCourseVideos(
	ctx context.Context,
	id int,
) (*coursemodel.CourseVideosResponse, error) {
	course, err := biz.store.FindOne(
		ctx,
		map[string]interface{}{"id": id},
		"Videos.Lesson",
		"Videos.ProcessInfos",
	)
	if err != nil {
		if err == common.RecordNotFound {
			return nil, common.ErrCannotGetEntity(models.CourseEntityName, err)
		}
		return nil, common.ErrCannotGetEntity(models.CourseEntityName, err)
	}

	var videos []coursemodel.CourseVideoResponse
	for _, video := range course.Videos {
		var lessonRes coursemodel.VideoLessonResonse
		copier.Copy(&lessonRes, video.Lesson)

		var processInfo []coursemodel.VideoProcessInfoResponse
		copier.Copy(&processInfo, video.ProcessInfos)

		videoRes := coursemodel.CourseVideoResponse{
			SQLModel:     video.SQLModel,
			Title:        video.Title,
			Description:  video.Description,
			ThumbnailURL: video.ThumbnailURL,
			Duration:     video.Duration,
			Order:        video.Order,
			AllowPreview: video.AllowPreview,
			Lesson:       lessonRes,
			ProcessInfos: processInfo,
		}
		videoRes.Mask(false)
		videos = append(videos, videoRes)
	}

	response := &coursemodel.CourseVideosResponse{
		Title:  course.Title,
		Videos: videos,
	}

	response.Mask(false)
	return response, nil
}
