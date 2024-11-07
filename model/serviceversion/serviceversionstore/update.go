package serviceversionstore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) Update(
	ctx context.Context,
	versionID uint32,
	updates *models.ServiceVersion,
) error {
	if err := s.db.Model(&models.ServiceVersion{}).
		Where("id = ?", versionID).
		Updates(updates).Error; err != nil {
		return common.ErrDB(err)
	}

	if updates.Images != nil {
		if err := s.db.
			Model(&models.ServiceVersion{}).
			Association("Images").
			Replace(updates.Images); err != nil {
			return common.ErrDB(err)
		}
	}

	return nil
}
