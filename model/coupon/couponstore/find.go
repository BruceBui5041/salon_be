package couponstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"

	"gorm.io/gorm"
)

func (s *sqlStore) Find(
	ctx context.Context,
	conditions map[string]interface{},
	moreInfo ...string,
) ([]models.Coupon, error) {
	var coupons []models.Coupon
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).Find(&coupons).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return coupons, nil
}
