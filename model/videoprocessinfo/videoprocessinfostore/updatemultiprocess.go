package videoprocessinfostore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) UpdateMultiProcessState(
	ctx context.Context,
	processInfos []*models.VideoProcessInfo,
) error {
	for _, processInfo := range processInfos {
		if err := s.db.Model(&models.VideoProcessInfo{}).
			Where("id = ?", processInfo.Id).
			Updates(processInfo).Error; err != nil {
			return err
		}
	}
	return nil
}
