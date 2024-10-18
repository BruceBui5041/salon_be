package coursestore

import (
	"context"
	"errors"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) CreateNewCourse(
	ctx context.Context,
	newCourse *models.Course,
) (uint32, error) {
	// Check if course with the same slug already exists
	var existingCourse models.Course
	if err := s.db.Where("slug = ?", newCourse.Slug).First(&existingCourse).Error; err == nil {
		return 0, errors.New("course with this slug already exists")
	} else if err != gorm.ErrRecordNotFound {
		return 0, err
	}

	// Create the new course
	if err := s.db.Create(newCourse).Error; err != nil {
		return 0, err
	}

	return newCourse.Id, nil
}
