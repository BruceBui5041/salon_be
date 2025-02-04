package ekycbiz

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/model/ekyc/ekycmodel"
)

type KYCProfileRepo interface {
	ProcessKYCProfile(ctx context.Context, input *ekycmodel.CreateKYCProfileRequest) error
}

type createKYCBiz struct {
	repo KYCProfileRepo
}

func NewCreateKYCBiz(repo KYCProfileRepo) *createKYCBiz {
	return &createKYCBiz{repo: repo}
}

func (biz *createKYCBiz) CreateKYCProfile(ctx context.Context, input *ekycmodel.CreateKYCProfileRequest) error {
	if input.UserID == 0 {
		return common.ErrInvalidRequest(errors.New("user ID is required"))
	}

	if input.FrontDocument == nil || input.BackDocument == nil || input.FaceImage == nil {
		return common.ErrInvalidRequest(errors.New("all document images are required"))
	}

	if input.ClientSession == "" {
		return common.ErrInvalidRequest(errors.New("client session is required"))
	}

	if err := biz.repo.ProcessKYCProfile(ctx, input); err != nil {
		return common.ErrCannotCreateEntity("KYCProfile", err)
	}

	return nil
}
