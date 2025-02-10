package locationbiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

type UpdateLocationRepository interface {
	FindLocationByUserID(ctx context.Context, userID uint32) (*models.Location, error)
	CreateLocation(ctx context.Context, data *models.Location) error
	UpdateLocation(ctx context.Context, data *models.Location) error
}

type updateLocationBiz struct {
	repo UpdateLocationRepository
}

func NewUpdateLocationBiz(repo UpdateLocationRepository) *updateLocationBiz {
	return &updateLocationBiz{repo: repo}
}

func (biz *updateLocationBiz) UpdateLocation(ctx context.Context, data *models.Location) error {
	existingLocation, err := biz.repo.FindLocationByUserID(ctx, data.UserId)
	if err != nil && err != common.RecordNotFound {
		return common.ErrCannotGetEntity(models.LocationEntityName, err)
	}

	if existingLocation == nil {
		if err := biz.repo.CreateLocation(ctx, data); err != nil {
			return common.ErrCannotCreateEntity(models.LocationEntityName, err)
		}
		return nil
	}

	if err := biz.repo.UpdateLocation(ctx, data); err != nil {
		return common.ErrCannotUpdateEntity(models.LocationEntityName, err)
	}

	return nil
}
