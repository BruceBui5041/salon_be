package videoprocessinfostore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) Create(
	ctx context.Context,
	processInfo *models.VideoProcessInfo,
) (uint32, error) {
	if err := s.db.Create(processInfo).Error; err != nil {
		return 0, err
	}

	return processInfo.Id, nil
}
