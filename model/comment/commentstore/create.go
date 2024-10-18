package commentstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) CreateNewComment(
	ctx context.Context,
	newComment *models.Comment,
) (uint32, error) {
	if err := s.db.Create(newComment).Error; err != nil {
		return 0, err
	}

	return newComment.Id, nil
}
