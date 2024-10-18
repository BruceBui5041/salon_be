package rolestore

import (
	"context"
	models "video_server/model"
	"video_server/model/role/rolemodel"
)

func (s *sqlStore) Update(ctx context.Context, id uint32, data *rolemodel.UpdateRole) error {
	if err := s.db.Where("id = ?", id).Updates(&models.Role{
		Name:        data.Name,
		Code:        data.Code,
		Description: data.Description,
	}).Error; err != nil {
		return err
	}

	return nil
}
