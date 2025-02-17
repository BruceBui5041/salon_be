package groupproviderbiz

import (
	"context"
	"errors"
	"mime/multipart"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/groupprovider/groupprovidermodel"

	"gorm.io/gorm"
)

type CreateRepository interface {
	FindUserByID(ctx context.Context, userID uint32) (*models.User, error)
	FindGroupProviderByOwner(ctx context.Context, ownerID uint32) (*models.GroupProvider, error)
	Create(ctx context.Context, data *models.GroupProvider, images []*multipart.FileHeader, creatorID uint32) error
	UpdateServicesGroupProvider(ctx context.Context, ownerID uint32, groupProviderID uint32) error
}

type createBiz struct {
	repo CreateRepository
}

func NewCreateBiz(repo CreateRepository) *createBiz {
	return &createBiz{repo: repo}
}

func (biz *createBiz) CreateGroupProvider(ctx context.Context, data *groupprovidermodel.GroupProviderCreate) error {
	// Verify requester is admin
	requester, err := biz.repo.FindUserByID(ctx, data.RequesterID)
	if err != nil {
		return common.ErrEntityNotFound(models.UserEntityName, err)
	}

	if !requester.IsAdmin() {
		return common.ErrNoPermission(errors.New("user must be admin"))
	}

	ownerUID, err := common.FromBase58(data.OwnerStrID)
	if err != nil {
		return common.ErrInvalidRequest(err)
	}

	owner, err := biz.repo.FindUserByID(ctx, ownerUID.GetLocalID())
	if err != nil {
		return common.ErrEntityNotFound(models.UserEntityName, err)
	}

	// Check if user already has a group provider
	existingGroupProvider, err := biz.repo.FindGroupProviderByOwner(ctx, owner.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return common.ErrDB(err)
	}

	if existingGroupProvider != nil {
		return common.ErrInvalidRequest(errors.New("user already has a group provider"))
	}

	groupProvider := &models.GroupProvider{
		Name:        data.Name,
		Code:        data.Code,
		Description: data.Description,
		OwnerID:     owner.Id,
		CreatorID:   requester.Id,
		SQLModel:    common.SQLModel{Status: common.StatusInactive},
		Admins:      []*models.User{owner},
	}

	if err := biz.repo.Create(ctx, groupProvider, data.Images, requester.Id); err != nil {
		return common.ErrCannotCreateEntity(models.GroupProviderEntityName, err)
	}

	if err := biz.repo.UpdateServicesGroupProvider(ctx, requester.Id, groupProvider.Id); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
