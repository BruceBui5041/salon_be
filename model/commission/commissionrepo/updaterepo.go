package commissionrepo

import (
	"context"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	commissionmodel "salon_be/model/commission/comissionmodel"
	"salon_be/model/commission/commissionerror"

	"go.uber.org/zap"
)

type UpdateCommissionStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Commission, error)
	Update(ctx context.Context, id uint32, data *models.Commission) error
}

type updateCommissionRepo struct {
	store UpdateCommissionStore
}

func NewUpdateCommissionRepo(store UpdateCommissionStore) *updateCommissionRepo {
	return &updateCommissionRepo{store: store}
}

func (repo *updateCommissionRepo) UpdateCommission(ctx context.Context, id uint32, data *commissionmodel.UpdateCommission) error {
	existingCommission, err := repo.store.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to find commission in database",
			zap.Error(err),
			zap.Uint32("id", id))
		return err
	}

	// Check if commission is published
	if existingCommission.PublishedAt != nil {
		return commissionerror.ErrCommissionPublished()
	}

	var status string
	if data.Status == nil {
		status = common.StatusActive
	} else {
		status = common.StatusInactive
	}

	commission := &models.Commission{
		SQLModel:   common.SQLModel{Id: id, Status: status},
		Code:       data.Code,
		RoleID:     data.RoleID,
		Percentage: data.Percentage,
		MinAmount:  data.MinAmount,
		MaxAmount:  data.MaxAmount,
	}

	if err := repo.store.Update(ctx, id, commission); err != nil {
		logger.AppLogger.Error(ctx, "Failed to update commission in database",
			zap.Error(err),
			zap.Uint32("id", id),
			zap.String("code", data.Code))
		return err
	}

	return nil
}
