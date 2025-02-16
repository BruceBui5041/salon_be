package groupproviderrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type CreateGroupProviderStore interface {
	Create(ctx context.Context, data *models.GroupProvider) error
	FindOne(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*models.GroupProvider, error)
}

type CreateUserStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*models.User, error)
}

type CreateServiceStore interface {
	UpdateMany(ctx context.Context, condition map[string]interface{}, data *models.Service) error
}

type createRepo struct {
	store        CreateGroupProviderStore
	userStore    CreateUserStore
	serviceStore CreateServiceStore
}

func NewCreateRepo(
	store CreateGroupProviderStore,
	userStore CreateUserStore,
	serviceStore CreateServiceStore,
) *createRepo {
	return &createRepo{
		store:        store,
		userStore:    userStore,
		serviceStore: serviceStore,
	}
}

func (r *createRepo) FindUserByID(ctx context.Context, userID uint32) (*models.User, error) {
	user, err := r.userStore.FindOne(ctx, map[string]interface{}{"id": userID}, "Roles")
	if err != nil {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}
	return user, nil
}

func (r *createRepo) FindGroupProviderByOwner(ctx context.Context, ownerID uint32) (*models.GroupProvider, error) {
	groupProvider, err := r.store.FindOne(ctx, map[string]interface{}{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	return groupProvider, nil
}

func (r *createRepo) Create(ctx context.Context, data *models.GroupProvider) error {
	if err := r.store.Create(ctx, data); err != nil {
		return err
	}
	return nil
}

func (r *createRepo) UpdateServicesGroupProvider(ctx context.Context, ownerID uint32, groupProviderID uint32) error {
	condition := map[string]interface{}{"owner_id": ownerID, "group_provider_id": nil}
	service := &models.Service{
		GroupProviderID: &groupProviderID,
	}

	if err := r.serviceStore.UpdateMany(ctx, condition, service); err != nil {
		return err
	}
	return nil
}
