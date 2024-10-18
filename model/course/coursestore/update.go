package coursestore

import (
	"context"
	"errors"
	models "video_server/model"

	"gorm.io/gorm"
)

func (s *sqlStore) Update(
	ctx context.Context,
	id uint32,
	updateData *models.Course,
) error {
	if updateData.Slug != "" {
		var existingCourse models.Course
		if err := s.db.Where("slug = ? AND id != ?", updateData.Slug, id).First(&existingCourse).Error; err == nil {
			return errors.New("course with this slug already exists")
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
	}

	if err := s.db.Model(&models.Course{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		return err
	}

	return nil
}
