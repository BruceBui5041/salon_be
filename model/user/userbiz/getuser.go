package userbiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/user/usermodel"

	"github.com/jinzhu/copier"
)

type GetUserRepo interface {
	GetUser(ctx context.Context, id uint32) (*models.User, error)
}

type getUserBiz struct {
	repo GetUserRepo
}

func NewGetUserBiz(repo GetUserRepo) *getUserBiz {
	return &getUserBiz{repo: repo}
}

func (biz *getUserBiz) GetUserById(
	ctx context.Context,
	id uint32,
) (*usermodel.GetUserResponse, error) {
	user, err := biz.repo.GetUser(ctx, id)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.UserEntityName, err)
	}

	user.Mask(false)
	for _, role := range user.Roles {
		role.Mask(false)
	}

	for _, enrollment := range user.Enrollments {
		enrollment.Mask(false)
		enrollment.ServiceVersion.Mask(false)
	}

	var userResponse usermodel.GetUserResponse
	copier.Copy(&userResponse, user)

	return &userResponse, nil
}
