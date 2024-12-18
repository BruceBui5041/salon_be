package userprofilestore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) CreateNewUserProfile(
	ctx context.Context,
	newUserProfile *models.UserProfile,
) (uint32, error) {
	if err := s.db.Create(newUserProfile).Error; err != nil {
		return 0, err
	}

	return newUserProfile.Id, nil
}
