package lecturestore

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/lecture/lecturemodel"
)

func (s *sqlStore) Update(
	ctx context.Context,
	lectureId uint32,
	data *lecturemodel.UpdateLecture,
) error {
	if err := s.db.Table(models.Lecture{}.TableName()).Where("id = ?", lectureId).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
