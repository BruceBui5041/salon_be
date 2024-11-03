package m2mserviceversionimagestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Delete(ctx context.Context, conditions map[string]interface{}) error {
	if err := s.db.Where(conditions).Delete(&models.M2MServiceVersionImage{}).Error; err != nil {
		return err
	}
	return nil
}
