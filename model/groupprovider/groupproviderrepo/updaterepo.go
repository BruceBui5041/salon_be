package groupproviderrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type UpdateGroupProviderStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*models.GroupProvider, error)
	Update(ctx context.Context, id uint32, data *models.GroupProvider) error
}

type UpdateUserStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*models.User, error)
	Find(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) ([]*models.User, error)
}

type updateRepo struct {
	store     UpdateGroupProviderStore
	userStore UpdateUserStore
}

func NewUpdateRepo(
	store UpdateGroupProviderStore,
	userStore UpdateUserStore,
) *updateRepo {
	return &updateRepo{
		store:     store,
		userStore: userStore,
	}
}

func (r *updateRepo) FindUserByID(ctx context.Context, userID uint32) (*models.User, error) {
	user, err := r.userStore.FindOne(ctx, map[string]interface{}{"id": userID}, "Roles")
	if err != nil {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}
	return user, nil
}

func (r *updateRepo) FindGroupProviderByID(ctx context.Context, id uint32) (*models.GroupProvider, error) {
	groupProvider, err := r.store.FindOne(ctx, map[string]interface{}{"id": id}, "Owner", "Admins")
	if err != nil {
		return nil, err
	}
	return groupProvider, nil
}

func (r *updateRepo) Update(ctx context.Context, id uint32, data *models.GroupProvider) error {
	if err := r.store.Update(ctx, id, data); err != nil {
		return err
	}
	return nil
}

func (r *updateRepo) FindUsersByIDs(ctx context.Context, ids []uint32) ([]*models.User, error) {
	users, err := r.userStore.Find(ctx, map[string]interface{}{"id": ids}, "Roles")
	if err != nil {
		return nil, err
	}
	return users, nil
}
