package groupproviderbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/groupprovider/groupprovidermodel"
)

type UpdateRepository interface {
	FindUserByID(ctx context.Context, userID uint32) (*models.User, error)
	FindGroupProviderByID(ctx context.Context, id uint32) (*models.GroupProvider, error)
	Update(ctx context.Context, id uint32, data *models.GroupProvider) error
	FindUsersByIDs(ctx context.Context, ids []uint32) ([]*models.User, error)
}

type updateBiz struct {
	repo UpdateRepository
}

func NewUpdateBiz(repo UpdateRepository) *updateBiz {
	return &updateBiz{repo: repo}
}

func (biz *updateBiz) UpdateGroupProvider(ctx context.Context, id uint32, data *groupprovidermodel.GroupProviderUpdate) error {
	requester, err := biz.repo.FindUserByID(ctx, data.RequesterID)
	if err != nil {
		return common.ErrEntityNotFound(models.UserEntityName, err)
	}

	groupProvider, err := biz.repo.FindGroupProviderByID(ctx, id)
	if err != nil {
		return common.ErrEntityNotFound(models.GroupProviderEntityName, err)
	}

	// Check if requester has permission
	if !requester.IsAdmin() && !requester.IsGroupProviderAdmin() && groupProvider.OwnerID != requester.Id {
		return common.ErrNoPermission(errors.New("user must be admin, group admin or owner"))
	}

	// Handle owner change
	if data.OwnerStrID != "" {
		if !requester.IsAdmin() {
			return common.ErrNoPermission(errors.New("only admin can change owner"))
		}

		ownerUID, err := common.FromBase58(data.OwnerStrID)
		if err != nil {
			return common.ErrInvalidRequest(err)
		}

		newOwner, err := biz.repo.FindUserByID(ctx, ownerUID.GetLocalID())
		if err != nil {
			return common.ErrEntityNotFound(models.UserEntityName, err)
		}

		groupProvider.OwnerID = newOwner.Id
		groupProvider.Owner = newOwner
	}

	// Handle admin changes
	if len(data.AdminIDs) > 0 {
		adminIDs := make([]uint32, len(data.AdminIDs))
		for i, adminID := range data.AdminIDs {
			uid, err := common.FromBase58(adminID)
			if err != nil {
				return common.ErrInvalidRequest(err)
			}
			adminIDs[i] = uid.GetLocalID()
		}

		admins, err := biz.repo.FindUsersByIDs(ctx, adminIDs)
		if err != nil {
			return common.ErrEntityNotFound(models.UserEntityName, err)
		}

		groupProvider.Admins = admins
	}

	groupProvider.Name = data.Name
	groupProvider.Description = data.Description

	if err := biz.repo.Update(ctx, id, groupProvider); err != nil {
		return common.ErrCannotUpdateEntity(models.GroupProviderEntityName, err)
	}

	return nil
}
