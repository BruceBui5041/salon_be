package videoprocessinfostore

import (
	"context"
	models "video_server/model"
)

func (s *sqlStore) CreateMultiProcessState(
	ctx context.Context,
	processInfos []*models.VideoProcessInfo,
) ([]uint32, error) {
	var ids []uint32

	for _, processInfo := range processInfos {
		if err := s.db.Create(processInfo).Error; err != nil {
			return nil, err
		}
		ids = append(ids, processInfo.Id)
	}

	return ids, nil
}
