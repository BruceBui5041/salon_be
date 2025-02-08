package ekycbiz

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/model/ekyc/ekycmodel"
)

type uploadRepo interface {
	UploadKYCImage(ctx context.Context, userId uint32, input *ekycmodel.UploadRequest) (*ekycmodel.KYCImageUploadRes, error)
}

type uploadBiz struct {
	repo uploadRepo
}

func NewUploadBiz(repo uploadRepo) *uploadBiz {
	return &uploadBiz{repo: repo}
}

func (biz *uploadBiz) UploadImage(
	ctx context.Context,
	userId uint32,
	input *ekycmodel.UploadRequest,
) (*ekycmodel.KYCImageUploadRes, error) {
	if input.Image == nil {
		return nil, common.ErrInvalidRequest(errors.New("image is required"))
	}

	return biz.repo.UploadKYCImage(ctx, userId, input)
}
