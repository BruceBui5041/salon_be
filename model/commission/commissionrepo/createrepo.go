package commissionrepo

import (
	"context"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	commissionmodel "salon_be/model/commission/comissionmodel"

	"go.uber.org/zap"
)

type CreateCommissionStore interface {
	Create(ctx context.Context, data *models.Commission) error
}

type createCommissionRepo struct {
	store CreateCommissionStore
}

func NewCreateCommissionRepo(store CreateCommissionStore) *createCommissionRepo {
	return &createCommissionRepo{store: store}
}

func (repo *createCommissionRepo) CreateCommission(ctx context.Context, data *commissionmodel.CreateCommission) (uint32, error) {
	var status string
	if data.Status == nil {
		status = common.StatusActive
	} else {
		status = common.StatusInactive
	}

	commission := &models.Commission{
		SQLModel:   common.SQLModel{Status: status},
		Code:       data.Code,
		RoleID:     data.RoleID,
		Percentage: data.Percentage,
		MinAmount:  data.MinAmount,
		MaxAmount:  data.MaxAmount,
	}

	if err := repo.store.Create(ctx, commission); err != nil {
		logger.AppLogger.Error(ctx, "Failed to create commission in database",
			zap.Error(err),
			zap.String("code", data.Code))
		return 0, err
	}

	return commission.Id, nil
}
