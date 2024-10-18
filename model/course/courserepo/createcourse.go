package courserepo

import (
	"context"
	"errors"
	"fmt"
	"video_server/appconst"
	"video_server/common"
	models "video_server/model"
	"video_server/model/course/coursemodel"
	"video_server/storagehandler"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type CategoryStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
}

type CreateCourseUserStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.User, error)
}

type CreateCourseStore interface {
	CreateNewCourse(
		ctx context.Context,
		newCourse *models.Course,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
	Update(
		ctx context.Context,
		id uint32,
		updateData *models.Course,
	) error
}

type createCourseRepo struct {
	courseStore   CreateCourseStore
	categoryStore CategoryStore
	userStore     CreateCourseUserStore
	s3Client      *s3.S3
}

func NewCreateCourseRepo(
	courseStore CreateCourseStore,
	categoryStore CategoryStore,
	userStore CreateCourseUserStore,
	s3Client *s3.S3,
) *createCourseRepo {
	return &createCourseRepo{
		courseStore:   courseStore,
		categoryStore: categoryStore,
		userStore:     userStore,
		s3Client:      s3Client,
	}
}

func (repo *createCourseRepo) CreateNewCourse(ctx context.Context, input *coursemodel.CreateCourse) (*models.Course, error) {
	uid, err := common.FromBase58(input.CategoryID)
	if err != nil {
		return nil, err
	}

	category, err := repo.categoryStore.FindOne(ctx, map[string]interface{}{"id": uid.GetLocalID()})
	if err != nil {
		return nil, err
	}

	newCourse := &models.Course{
		CategoryID: category.Id,
	}

	if err := copier.Copy(newCourse, input); err != nil {
		return nil, err
	}

	courseId, err := repo.courseStore.CreateNewCourse(ctx, newCourse)
	if err != nil {
		return nil, err
	}

	course, err := repo.courseStore.FindOne(ctx, map[string]interface{}{"id": courseId}, "Creator")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	if input.Thumbnail != nil {
		thumbnailFile, err := input.Thumbnail.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open thumbnail file: %w", err)
		}

		defer thumbnailFile.Close()

		input.Mask(false)

		thumbnailKey := storagehandler.GenerateCourseThumbnaiS3Key(
			storagehandler.CourseInfo{
				UploadedBy: course.Creator.GetFakeId(),
				CourseId:   course.GetFakeId(),
				Filename:   input.Thumbnail.Filename,
			},
		)

		err = storagehandler.UploadFileToS3(ctx, repo.s3Client, thumbnailFile, appconst.AWSPublicBucket, thumbnailKey)
		if err != nil {
			return nil, fmt.Errorf("failed to upload thumbnail to S3: %w", err)
		}

		course.Thumbnail = thumbnailKey
	}

	if err := repo.courseStore.Update(ctx, courseId, course); err != nil {
		storagehandler.RemoveFileFromS3(ctx, repo.s3Client, appconst.AWSPublicBucket, course.Thumbnail)
		return nil, err
	}

	course, err = repo.courseStore.FindOne(ctx, map[string]interface{}{"id": courseId})
	if err != nil {
		return nil, err
	}

	return course, nil
}
