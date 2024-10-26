package imagerepo

import (
	"context"
	"mime/multipart"
	"salon_be/appconst"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/storagehandler"

	"github.com/aws/aws-sdk-go/service/s3"
)

type ImageStore interface {
	Create(ctx context.Context, data *models.Image) error
}

type createImageRepo struct {
	store    ImageStore
	s3Client *s3.S3
}

func NewCreateImageRepo(
	store ImageStore,
	s3Client *s3.S3,
) *createImageRepo {
	return &createImageRepo{
		store:    store,
		s3Client: s3Client,
	}
}

func (repo *createImageRepo) CreateImage(
	ctx context.Context,
	file *multipart.FileHeader,
	serviceID uint32,
	userID uint32,
) (*models.Image, error) {
	fileBytes, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileBytes.Close()

	tempObj := common.SQLModel{Id: serviceID}
	tempObj.GenUID(common.DBTypeService)

	objectKey := storagehandler.GenerateServiceImageS3Key(tempObj.GetFakeId(), file.Filename)
	if err := storagehandler.UploadFileToS3(ctx, repo.s3Client, fileBytes, appconst.AWSPublicBucket, objectKey); err != nil {
		return nil, err
	}

	img := &models.Image{
		UserID:    userID,
		ServiceID: serviceID,
		URL:       objectKey,
	}

	if err := repo.store.Create(ctx, img); err != nil {
		return nil, err
	}

	return img, nil
}
