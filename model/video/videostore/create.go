package videostore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) CreateNewVideo(
	ctx context.Context,
	newVideo *models.Video,
) (uint32, error) {
	if err := s.db.Create(newVideo).Error; err != nil {
		return 0, err
	}

	return newVideo.Id, nil
}
