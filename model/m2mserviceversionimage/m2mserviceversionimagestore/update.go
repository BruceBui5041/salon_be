package m2mserviceversionimagestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(ctx context.Context, updates *models.M2MServiceVersionImage) error {
	return s.db.Model(&models.M2MServiceVersionImage{}).
		Where("image_id = ? AND service_version_id = ?", updates.ImageID, updates.ServiceVersionID).
		Updates(updates).Error
}
