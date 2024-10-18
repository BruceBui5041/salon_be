package videobiz

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/appconst"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/video/videomodel"
	"salon_be/storagehandler"

	"github.com/aws/aws-sdk-go/service/s3"
)

type UpdateVideoRepo interface {
	UpdateVideo(ctx context.Context, id uint32, input *videomodel.UpdateVideo) (*models.Video, error)
	GetVideo(ctx context.Context, id uint32) (*models.Video, error)
	GetS3Client() *s3.S3
}

type updateVideoBiz struct {
	repo UpdateVideoRepo
}

func NewUpdateVideoBiz(repo UpdateVideoRepo) *updateVideoBiz {
	return &updateVideoBiz{repo: repo}
}

func (v *updateVideoBiz) UpdateVideo(
	ctx context.Context,
	id uint32,
	input *videomodel.UpdateVideo,
	videoFile *multipart.FileHeader,
	thumbnailFile *multipart.FileHeader,
	userId uint32,
) (*models.Video, error) {
	// Validate input
	if err := v.validateInput(input); err != nil {
		return nil, err
	}

	// Get existing video
	existingVideo, err := v.repo.GetVideo(ctx, id)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.VideoEntityName, err)
	}

	// Handle file uploads
	if err := v.handleFileUploads(
		ctx,
		input,
		videoFile,
		thumbnailFile,
		userId,
		existingVideo.Id,
	); err != nil {
		return nil, err
	}

	// Update video
	video, err := v.repo.UpdateVideo(ctx, id, input)
	if err != nil {
		return nil, common.ErrCannotUpdateEntity(models.VideoEntityName, err)
	}

	return video, nil
}

func (v *updateVideoBiz) validateInput(input *videomodel.UpdateVideo) error {
	if input.Title != nil && *input.Title == "" {
		return errors.New("video title cannot be empty")
	}

	if input.Title != nil && len(*input.Title) > 255 {
		return errors.New("video title must not exceed 255 characters")
	}

	return nil
}

func (v *updateVideoBiz) handleFileUploads(
	ctx context.Context,
	input *videomodel.UpdateVideo,
	videoFile *multipart.FileHeader,
	thumbnailFile *multipart.FileHeader,
	userId uint32,
	videoId uint32,
) error {
	fakeVideo := common.SQLModel{Id: videoId}
	fakeVideo.GenUID(common.DbTypeVideo)

	fakeUsr := common.SQLModel{Id: userId}
	fakeUsr.GenUID(common.DbTypeUser)

	if videoFile != nil {
		var videoReader interface{ Read([]byte) (int, error) }
		videoReader, err := videoFile.Open()
		if err != nil {
			return err
		}
		defer videoReader.(interface{ Close() error }).Close()

		videoStorageInfo := storagehandler.VideoInfo{
			UploadedBy:    fakeUsr.FakeId.String(),
			VideoId:       fakeVideo.FakeId.String(),
			VideoFilename: videoFile.Filename,
		}

		videoKey := storagehandler.GenerateVideoS3Key(videoStorageInfo)
		err = storagehandler.UploadFileToS3(
			ctx,
			v.repo.GetS3Client(),
			videoReader,
			appconst.AWSVideoS3BuckerName,
			videoKey,
		)
		if err != nil {
			return errors.New("failed to upload video to S3")
		}

		input.VideoURL = &videoKey
	}

	if thumbnailFile != nil {
		var thumbnailReader interface{ Read([]byte) (int, error) }
		thumbnailReader, err := thumbnailFile.Open()
		if err != nil {
			return err
		}
		defer thumbnailReader.(interface{ Close() error }).Close()

		if thumbnailReader != nil {
			thumbnailStorageInfo := storagehandler.VideoInfo{
				UploadedBy:        fakeUsr.FakeId.String(),
				VideoId:           fakeVideo.FakeId.String(),
				ThumbnailFilename: thumbnailFile.Filename,
			}

			thumbnailKey := storagehandler.GenerateVideoThumbnailS3Key(thumbnailStorageInfo)
			err := storagehandler.UploadFileToS3(
				ctx,
				v.repo.GetS3Client(),
				thumbnailReader,
				appconst.AWSVideoS3BuckerName,
				thumbnailKey,
			)
			if err != nil {
				return errors.New("failed to upload thumbnail to S3")
			}

			input.ThumbnailURL = &thumbnailKey
		}
	}

	return nil
}
