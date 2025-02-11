package locationbiz

import (
	"context"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"

	"go.uber.org/zap"
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
	logger.AppLogger.Info(
		ctx,
		"location_update",
		zap.Any("input", data),
		zap.Any("existing", existingLocation),
	)
	if err != nil && err != common.RecordNotFound {
		if existingLocation == nil {
			if err := biz.repo.CreateLocation(ctx, data); err != nil {
				return common.ErrCannotCreateEntity(models.LocationEntityName, err)
			}
			return nil
		}

		return common.ErrCannotGetEntity(models.LocationEntityName, err)
	}

	data.Id = existingLocation.Id
	if err := biz.repo.UpdateLocation(ctx, data); err != nil {
		return common.ErrCannotUpdateEntity(models.LocationEntityName, err)
	}

	return nil
}
