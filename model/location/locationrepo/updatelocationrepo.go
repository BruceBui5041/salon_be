package locationrepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type UpdateLocationStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreKeys ...string,
	) (*models.Location, error)
	Create(ctx context.Context, data *models.Location) error
	Update(ctx context.Context, data *models.Location) error
}

type updateLocationRepo struct {
	store UpdateLocationStore
}

func NewUpdateLocationRepo(store UpdateLocationStore) *updateLocationRepo {
	return &updateLocationRepo{store: store}
}

func (r *updateLocationRepo) FindLocationByUserID(ctx context.Context, userID uint32) (*models.Location, error) {
	conditions := map[string]interface{}{"user_id": userID}
	location, err := r.store.FindOne(ctx, conditions)
	if err != nil {
		if err == common.RecordNotFound {
			return nil, err
		}
		return nil, common.ErrDB(err)
	}
	return location, nil
}

func (r *updateLocationRepo) CreateLocation(ctx context.Context, data *models.Location) error {
	if err := r.store.Create(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (r *updateLocationRepo) UpdateLocation(ctx context.Context, data *models.Location) error {
	if err := r.store.Update(ctx, data); err != nil {
		return common.ErrDB(err)
	}
	return nil
}
