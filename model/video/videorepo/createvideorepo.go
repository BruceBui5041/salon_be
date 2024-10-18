package videorepo

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/video/videomodel"
	"salon_be/storagehandler"
	"salon_be/utils"

	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type CreateVideoCourseStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...interface{},
	) (*models.Course, error)
}

type CreateVideoStore interface {
	CreateNewVideo(
		ctx context.Context,
		newVideo *models.Video,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Video, error)
	UpdateVideo(
		ctx context.Context,
		id uint32,
		updateData *videomodel.UpdateVideo,
	) error
}

type CreateVideoProcessStore interface {
	CreateMultiProcessState(
		ctx context.Context,
		processInfos []*models.VideoProcessInfo,
	) ([]uint32, error)
}

type createVideoRepo struct {
	videoStore        CreateVideoStore
	courseStore       CreateVideoCourseStore
	videoProcessStore CreateVideoProcessStore
	svc               *s3.S3
}

func NewCreateVideoRepo(
	videoStore CreateVideoStore,
	courseStore CreateVideoCourseStore,
	videoProcessStore CreateVideoProcessStore,
	svc *s3.S3,
) *createVideoRepo {
	return &createVideoRepo{
		videoStore:        videoStore,
		courseStore:       courseStore,
		videoProcessStore: videoProcessStore,
		svc:               svc,
	}
}

func (repo *createVideoRepo) CreateNewVideo(
	ctx context.Context,
	input *videomodel.CreateVideo,
	videoFile,
	thumbnailFile *multipart.FileHeader,
) (*models.Video, error) {
	uid, err := common.FromBase58(input.CourseID)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to parse CourseID", zap.Error(err))
		return nil, err
	}

	course, err := repo.courseStore.FindOne(ctx, map[string]interface{}{"id": uid.GetLocalID()})
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to find course", zap.Error(err))
		return nil, err
	}

	newVideo := &models.Video{
		CourseID:     course.Id,
		Title:        input.Title,
		Description:  input.Description,
		VideoURL:     input.VideoURL,
		Duration:     input.Duration,
		Order:        input.Order,
		ThumbnailURL: input.ThumbnailURL,
	}

	videoId, err := repo.videoStore.CreateNewVideo(ctx, newVideo)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to create new video", zap.Error(err))
		return nil, err
	}

	video, err := repo.videoStore.FindOne(ctx, map[string]interface{}{"id": videoId})
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to find created video", zap.Error(err))
		return nil, err
	}

	video.Mask(false)
	course.Mask(false)

	sqlObj := common.SQLModel{Id: course.CreatorID}
	sqlObj.GenUID(common.DbTypeUser)

	videoStorageInfo := storagehandler.VideoInfo{
		UploadedBy:        sqlObj.FakeId.String(),
		CourseId:          course.FakeId.String(),
		VideoId:           video.FakeId.String(),
		ThumbnailFilename: thumbnailFile.Filename,
		VideoFilename:     videoFile.Filename,
	}

	videoKey := storagehandler.GenerateVideoS3Key(videoStorageInfo)
	thumbnailKey := storagehandler.GenerateVideoThumbnailS3Key(videoStorageInfo)

	if err := repo.uploadFiles(ctx, videoFile, thumbnailFile, videoKey, thumbnailKey); err != nil {
		logger.AppLogger.Error(ctx, "failed to upload files", zap.Error(err))
		repo.removeFiles(ctx, videoKey, thumbnailKey)
		return nil, err
	}

	video.RawVideoURL = videoKey
	video.ThumbnailURL = thumbnailKey
	video.VideoURL = utils.RemoveFileExtension(videoKey)
	err = repo.videoStore.UpdateVideo(
		ctx,
		videoId,
		&videomodel.UpdateVideo{
			RawVideoURL:  videoKey,
			VideoURL:     &video.VideoURL,
			ThumbnailURL: &thumbnailKey,
		},
	)

	if err != nil {
		logger.AppLogger.Error(ctx, "failed to update video", zap.Error(err))
		repo.removeFiles(ctx, videoKey, thumbnailKey)
		return nil, err
	}

	processStates := []*models.VideoProcessInfo{
		{
			VideoID:           videoId,
			ProcessResolution: "360p",
		},
		{
			VideoID:           videoId,
			ProcessResolution: "480p",
		},
		{
			VideoID:           videoId,
			ProcessResolution: "720p",
		},
		{
			VideoID:           videoId,
			ProcessResolution: "1080p",
		},
	}
	_, err = repo.videoProcessStore.CreateMultiProcessState(ctx, processStates)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to create process states", zap.Error(err))
		return nil, err
	}

	video.Course = *course

	return video, nil
}

func (repo *createVideoRepo) uploadFiles(ctx context.Context, videoFile, thumbnailFile *multipart.FileHeader, videoKey, thumbnailKey string) error {
	videoFileContent, err := videoFile.Open()
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to open video file", zap.Error(err))
		return errors.New("failed to open video file")
	}
	defer videoFileContent.Close()

	err = storagehandler.UploadFileToS3(
		ctx,
		repo.svc,
		videoFileContent,
		appconst.AWSVideoS3BuckerName,
		videoKey,
	)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to upload video to S3", zap.Error(err))
		return errors.New("failed to upload video to S3")
	}

	thumbnailFileContent, err := thumbnailFile.Open()
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to open thumbnail file", zap.Error(err))
		repo.removeFiles(ctx, videoKey, "")
		return errors.New("failed to open thumbnail file")
	}
	defer thumbnailFileContent.Close()

	err = storagehandler.UploadFileToS3(
		ctx,
		repo.svc,
		thumbnailFileContent,
		appconst.AWSPublicBucket,
		thumbnailKey,
	)
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to upload thumbnail to S3", zap.Error(err))
		repo.removeFiles(ctx, "", thumbnailKey)
		return errors.New("failed to upload thumbnail to S3")
	}

	return nil
}

func (repo *createVideoRepo) removeFiles(ctx context.Context, videoKey, thumbnailKey string) {
	if videoKey != "" {
		err := storagehandler.RemoveFileFromS3(ctx, repo.svc, appconst.AWSVideoS3BuckerName, videoKey)
		if err != nil {
			logger.AppLogger.Error(ctx, "failed to remove video file from S3", zap.Error(err))
		}
	}
	if thumbnailKey != "" {
		err := storagehandler.RemoveFileFromS3(ctx, repo.svc, appconst.AWSPublicBucket, thumbnailKey)
		if err != nil {
			logger.AppLogger.Error(ctx, "failed to remove thumbnail file from S3", zap.Error(err))
		}
	}
}
