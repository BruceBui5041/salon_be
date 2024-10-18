package userprofilestore

import (
	"context"
	models "video_server/model"

	"github.com/jinzhu/copier"
)

func (s *sqlStore) UpdateProfile(
	ctx context.Context,
	profileId uint32,
	data *models.UserProfile,
) error {
	if err := copier.Copy(&data, data); err != nil {
		return err
	}

	if err := s.db.Where("id = ?", profileId).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
