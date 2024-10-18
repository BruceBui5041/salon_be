package userbiz

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/user/usermodel"

	"gorm.io/gorm"
)

type UpdateUserRepo interface {
	UpdateUser(ctx context.Context, tx *gorm.DB, id uint32, data *usermodel.UserUpdate) error
}

type updateUserBiz struct {
	repo UpdateUserRepo
}

func NewUpdateUserBiz(repo UpdateUserRepo) *updateUserBiz {
	return &updateUserBiz{repo: repo}
}

func (biz *updateUserBiz) UpdateUser(ctx context.Context, tx *gorm.DB, id uint32, data *usermodel.UserUpdate) error {
	if err := biz.repo.UpdateUser(ctx, tx, id, data); err != nil {
		return common.ErrCannotUpdateEntity(models.UserEntityName, err)
	}
	return nil
}
