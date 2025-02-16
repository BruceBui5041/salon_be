package imagerepo

import (
	"context"
	"mime/multipart"
	models "salon_be/model"

	"github.com/aws/aws-sdk-go/service/s3"
)

type UpdateImageStore interface {
	Create(ctx context.Context, data *models.Image) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.Image, error)
	List(
		ctx context.Context,
		conditions []interface{},
		moreKeys ...string,
	) ([]*models.Image, error)
}

type updateImageRepo struct {
	store    UpdateImageStore
	s3Client *s3.S3
}

func NewUpdateImageRepo(
	store UpdateImageStore,
	s3Client *s3.S3,
) *updateImageRepo {
	return &updateImageRepo{
		store:    store,
		s3Client: s3Client,
	}
}

func (repo *updateImageRepo) UpdateServiceImage(
	ctx context.Context,
	file *multipart.FileHeader,
	serviceID uint32,
	userID uint32,
) (*models.Image, error) {
	createImageRepo := NewCreateImageRepo(repo.store, repo.s3Client)

	image, err := createImageRepo.CreateServiceImage(ctx, file, serviceID, userID)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (repo *updateImageRepo) UpdateCouponImage(
	ctx context.Context,
	file *multipart.FileHeader,
	couponId uint32,
	userID uint32,
) (*models.Image, error) {
	createImageRepo := NewCreateImageRepo(repo.store, repo.s3Client)
	image, err := createImageRepo.CreateImageForCoupon(ctx, file, couponId, userID)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (repo *updateImageRepo) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreKeys ...string,
) (*models.Image, error) {
	return repo.store.FindOne(ctx, conditions, moreKeys...)
}

func (repo *updateImageRepo) List(
	ctx context.Context,
	conditions []interface{},
	moreKeys ...string,
) ([]*models.Image, error) {
	return repo.store.List(ctx, conditions, moreKeys...)
}
