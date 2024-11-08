package otpstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(ctx context.Context, updates *models.OTP) error {
	return s.db.Model(&models.OTP{}).
		Where("id = ?", updates.Id).
		Updates(updates).Error
}
