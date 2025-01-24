package commissionbiz

import (
	"context"
	commissionmodel "salon_be/model/commission/comissionmodel"
)

type UpdateCommissionRepo interface {
	UpdateCommission(ctx context.Context, id uint32, data *commissionmodel.UpdateCommission) error
}

type updateCommissionBiz struct {
	repo UpdateCommissionRepo
}

func NewUpdateCommissionBiz(repo UpdateCommissionRepo) *updateCommissionBiz {
	return &updateCommissionBiz{repo: repo}
}

func (biz *updateCommissionBiz) UpdateCommission(ctx context.Context, id uint32, data *commissionmodel.UpdateCommission) error {
	return biz.repo.UpdateCommission(ctx, id, data)
}
