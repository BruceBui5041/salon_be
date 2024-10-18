package lessonstore

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

func (s *sqlStore) Delete(ctx context.Context, lessonId uint32) error {
	if err := s.db.Table(models.Lesson{}.TableName()).Where("id = ?", lessonId).Delete(nil).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
