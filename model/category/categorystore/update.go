package categorystore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Update(
	ctx context.Context,
	id uint32,
	data *models.Category,
) error {
	if err := s.db.Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (s *sqlStore) UpdateParentId(
	ctx context.Context,
	id uint32,
	parentID *uint32,
) error {
	if err := s.db.
		Table(models.Category{}.TableName()).
		Where("id = ?", id).
		Updates(map[string]interface{}{"parent_id": parentID}).
		Error; err != nil {
		return err
	}

	return nil
}
