package rolebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/role/rolemodel"
)

type RoleRepo interface {
	CreateNewRole(ctx context.Context, input *rolemodel.CreateRole) (*models.Role, error)
}

type createRoleBiz struct {
	repo RoleRepo
}

func NewCreateRoleBiz(repo RoleRepo) *createRoleBiz {
	return &createRoleBiz{repo: repo}
}

func (c *createRoleBiz) CreateNewRole(ctx context.Context, input *rolemodel.CreateRole) error {
	if input.Name == "" {
		return common.ErrInvalidRequest(errors.New("role name is required"))
	}

	if len(input.Name) > 50 {
		return common.ErrInvalidRequest(errors.New("role name must not exceed 50 characters"))
	}

	if input.Code == "" {
		return common.ErrInvalidRequest(errors.New("role code is required"))
	}

	if len(input.Code) > 50 {
		return common.ErrInvalidRequest(errors.New("role code must not exceed 50 characters"))
	}

	_, err := c.repo.CreateNewRole(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.RoleEntityName, err)
	}

	return nil
}
