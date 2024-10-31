package rolestore

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
)

func (s *sqlStore) SoftDelete(ctx context.Context, id uint32) error {
	if err := s.db.Table(models.Role{}.TableName()).Where("id = ?", id).Update("status", common.StatusInactive).Error; err != nil {
		return err
	}

	return nil
}
