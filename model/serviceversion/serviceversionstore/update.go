package serviceversionstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(
	ctx context.Context,
	versionID uint32,
	updates *models.ServiceVersion,
) error {
	return s.db.Model(&models.ServiceVersion{}).
		Where("id = ?", versionID).
		Updates(updates).Error
}
