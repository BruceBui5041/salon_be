package bookingstore

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
) ([]models.Booking, error) {
	var bookings []models.Booking
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).Find(&bookings).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.RecordNotFound
		}

		return nil, common.ErrDB(err)
	}

	return bookings, nil
}
