package commentstore

import (
	"context"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/modelsearch"
	models "video_server/model"
)

func (s *sqlStore) Search(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
) ([]*models.Comment, error) {
	var comments []*models.Comment
	db := s.db

	query := modelsearch.Search(
		ctx,
		db.Model(&models.Comment{}),
		input.Model,
		input.Conditions,
		input.Fields,
		input.OrderBy,
	)

	err := query.Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}
