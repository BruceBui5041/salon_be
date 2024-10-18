package videorepo

import (
	"context"
	models "salon_be/model"
	"salon_be/model/video/videomodel"

	"github.com/aws/aws-sdk-go/service/s3"
)

type UpdateVideoStore interface {
	UpdateVideo(
		ctx context.Context,
		id uint32,
		updateData *videomodel.UpdateVideo,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Video, error)
}

type updateVideoRepo struct {
	store UpdateVideoStore
	svc   *s3.S3
}

func NewUpdateVideoRepo(store UpdateVideoStore, svc *s3.S3) *updateVideoRepo {
	return &updateVideoRepo{store: store, svc: svc}
}

func (repo *updateVideoRepo) GetS3Client() *s3.S3 {
	return repo.svc
}

func (repo *updateVideoRepo) UpdateVideo(ctx context.Context, id uint32, input *videomodel.UpdateVideo) (*models.Video, error) {
	if err := repo.store.UpdateVideo(ctx, id, input); err != nil {
		return nil, err
	}

	return repo.GetVideo(ctx, id)
}

func (repo *updateVideoRepo) GetVideo(ctx context.Context, id uint32) (*models.Video, error) {
	return repo.store.FindOne(ctx, map[string]interface{}{"id": id})
}
