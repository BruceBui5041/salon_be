package coursestore

import (
	"context"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/modelsearch"
	models "video_server/model"
)

func (s *sqlStore) Search(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
) ([]*models.Course, error) {
	var courses []*models.Course
	db := s.db

	query := modelsearch.Search(
		ctx,
		db.Model(&models.Course{}),
		input.Model,
		input.Conditions,
		input.Fields,
		input.OrderBy,
	)

	err := query.Find(&courses).Error
	if err != nil {
		return nil, err
	}

	return courses, nil
}
