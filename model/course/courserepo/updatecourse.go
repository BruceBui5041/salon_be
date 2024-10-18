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
)

type UpdateCourseVideoStore interface {
	Exist(
		ctx context.Context,
		conditions map[string]interface{},
	) (bool, error)
}

type UpdateCategoryStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
}

type UpdateCourseStore interface {
	Update(
		ctx context.Context,
		id uint32,
		updateData *models.Course,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type updateCourseRepo struct {
	updateCourseStore UpdateCourseStore
	categoryStore     UpdateCategoryStore
	videoStore        UpdateCourseVideoStore
	s3Client          *s3.S3
}

func NewUpdateCourseRepo(
	updateCourseStore UpdateCourseStore,
	categoryStore UpdateCategoryStore,
	videoStore UpdateCourseVideoStore,
	s3Client *s3.S3,
) *updateCourseRepo {
	return &updateCourseRepo{
		updateCourseStore: updateCourseStore,
		categoryStore:     categoryStore,
		videoStore:        videoStore,
		s3Client:          s3Client,
	}
}

func (repo *updateCourseRepo) UpdateCourse(ctx context.Context, courseId uint32, input *coursemodel.UpdateCourse) error {
	course, err := repo.updateCourseStore.FindOne(ctx, map[string]interface{}{"id": courseId}, "Creator")
	if err != nil {
		return err
	}

	uid, err := common.FromBase58(input.UploadedBy)
	if err != nil {
		return common.ErrInternal(err)
	}

	requesterId := uid.GetLocalID()
	if course.Creator.Id != requesterId {
		return common.ErrNoPermission(errors.New("only author can update the course"))
	}

	updateData := &models.Course{
		SQLModel:        common.SQLModel{Status: input.Status},
		Title:           input.Title,
		Description:     input.Description,
		Price:           input.Price.GetDecimal(),
		DiscountedPrice: input.DiscountedPrice.GetDecimal(),
		DifficultyLevel: input.DifficultyLevel,
		Overview:        input.Overview,
	}

	if input.IntroVideoId != "" {
		uid, err := common.FromBase58(input.IntroVideoId)
		if err != nil {
			return common.ErrInternal(err)
		}

		videoId := uid.GetLocalID()
		isExist, err := repo.videoStore.Exist(ctx, map[string]interface{}{"id": videoId, "course_id": courseId})
		if err != nil {
			return common.ErrDB(err)
		}

		if !isExist {
			return common.ErrEntityNotFound(models.VideoEntityName, errors.New("intro video not found"))
		}

		updateData.IntroVideoID = &videoId
	}

	if input.CategoryID != "" {
		uid, err := common.FromBase58(input.CategoryID)
		if err != nil {
			return err
		}
		category, err := repo.categoryStore.FindOne(ctx, map[string]interface{}{"id": uid.GetLocalID()})
		if err != nil {
			return err
		}
		updateData.CategoryID = category.Id
	}

	if input.Thumbnail != nil {
		thumbnailFile, err := input.Thumbnail.Open()
		if err != nil {
			return fmt.Errorf("failed to open thumbnail file: %w", err)
		}
		defer thumbnailFile.Close()

		input.Mask(false)

		thumbnailKey := storagehandler.GenerateCourseThumbnaiS3Key(
			storagehandler.CourseInfo{
				UploadedBy: input.UploadedBy,
				CourseId:   input.GetFakeId(),
				Filename:   input.Thumbnail.Filename,
			},
		)

		err = storagehandler.UploadFileToS3(ctx, repo.s3Client, thumbnailFile, appconst.AWSPublicBucket, thumbnailKey)
		if err != nil {
			return fmt.Errorf("failed to upload thumbnail to S3: %w", err)
		}

		updateData.Thumbnail = thumbnailKey
	}

	if err := repo.updateCourseStore.Update(ctx, courseId, updateData); err != nil {
		storagehandler.RemoveFileFromS3(ctx, repo.s3Client, appconst.AWSPublicBucket, updateData.Thumbnail)
		return err
	}

	return nil
}
