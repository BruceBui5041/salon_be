package enrollmentstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	newEnrollment *models.Enrollment,
) (uint32, error) {
	if err := s.db.Create(newEnrollment).Error; err != nil {
		return 0, err
	}

	return newEnrollment.Id, nil
}
