package userrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type GetUserStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*models.User, error)
}

type getUserRepo struct {
	store GetUserStore
}

func NewGetUserRepo(store GetUserStore) *getUserRepo {
	return &getUserRepo{
		store: store,
	}
}

func (repo *getUserRepo) GetUser(ctx context.Context, id uint32) (*models.User, error) {
	user, err := repo.store.FindOne(
		ctx,
		map[string]interface{}{"id": id},
		"Roles",
		"Enrollments.Service.Creator",
		"Enrollments.Service.Category",
		"UserProfile",
	)
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.UserEntityName, err)
	}

	for _, role := range user.Roles {
		role.Mask(false)
	}

	for _, enrollment := range user.Enrollments {
		enrollment.Mask(false)
	}

	if user.UserProfile != nil {
		user.UserProfile.Mask(false)
	}

	return user, nil
}
