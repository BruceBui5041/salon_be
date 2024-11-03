package servicebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/service/servicemodel"
)

type UploadImagesRepo interface {
	UploadImages(ctx context.Context, data *servicemodel.UploadImages) error
}

type uploadImagesBiz struct {
	repo UploadImagesRepo
}

func NewUploadImagesBiz(repo UploadImagesRepo) *uploadImagesBiz {
	return &uploadImagesBiz{repo: repo}
}

func (biz *uploadImagesBiz) UploadImages(ctx context.Context, data *servicemodel.UploadImages) error {
	if data.UploadedBy == 0 {
		return common.ErrInvalidRequest(errors.New("uploader is required"))
	}

	if data.ServiceID == "" {
		return common.ErrInvalidRequest(errors.New("service ID is required"))
	}

	if len(data.Images) == 0 {
		return common.ErrInvalidRequest(errors.New("images are required"))
	}

	if err := biz.repo.UploadImages(ctx, data); err != nil {
		return common.ErrCannotCreateEntity(models.ImageEntityName, err)
	}

	return nil
}
