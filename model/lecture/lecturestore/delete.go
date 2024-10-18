package lecturestore

import (
	"context"
	"video_server/common"
	models "video_server/model"
)

func (s *sqlStore) Delete(ctx context.Context, lectureId uint32) error {
	if err := s.db.Table(models.Lecture{}.TableName()).Where("id = ?", lectureId).Delete(nil).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
