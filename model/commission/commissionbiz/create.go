package commissionbiz

import (
	"context"
	commissionmodel "salon_be/model/commission/comissionmodel"
)

type CreateCommissionRepo interface {
	CreateCommission(ctx context.Context, data *commissionmodel.CreateCommission) (uint32, error)
}

type createCommissionBiz struct {
	repo CreateCommissionRepo
}

func NewCreateCommissionBiz(repo CreateCommissionRepo) *createCommissionBiz {
	return &createCommissionBiz{repo: repo}
}

func (biz *createCommissionBiz) CreateCommission(ctx context.Context, data *commissionmodel.CreateCommission) (uint32, error) {
	return biz.repo.CreateCommission(ctx, data)
}
