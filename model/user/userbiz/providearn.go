package userbiz

import (
	"context"
	"errors"
	"salon_be/common"
	"salon_be/model/user/usermodel"
	"time"
)

type ProviderEarningsRepo interface {
	CalculateEarnings(
		ctx context.Context,
		providerID uint32,
		fromDate, toDate time.Time,
	) (*usermodel.EarningsSummary, error)
}

type providerEarningsBiz struct {
	repo ProviderEarningsRepo
}

func NewProviderEarningsBiz(repo ProviderEarningsRepo) *providerEarningsBiz {
	return &providerEarningsBiz{repo: repo}
}

func (biz *providerEarningsBiz) GetProviderEarnings(
	ctx context.Context,
	providerID uint32,
	fromDate, toDate time.Time,
) (*usermodel.EarningsSummary, error) {
	if !fromDate.Before(toDate) {
		return nil, common.ErrInvalidRequest(errors.New("invalid date range"))
	}

	summary, err := biz.repo.CalculateEarnings(ctx, providerID, fromDate, toDate)
	if err != nil {
		return nil, common.ErrCannotListEntity("provider earnings", err)
	}

	return summary, nil
}
