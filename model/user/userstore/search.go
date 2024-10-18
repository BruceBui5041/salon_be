package userstore

import (
	"context"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/modelsearch"
	models "video_server/model"
)

func (s *sqlStore) Search(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
) ([]*models.User, error) {
	var users []*models.User
	db := s.db

	query := modelsearch.Search(
		ctx,
		db.Model(&models.User{}),
		input.Model,
		input.Conditions,
		input.Fields,
		input.OrderBy,
	)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
