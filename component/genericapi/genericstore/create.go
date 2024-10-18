package genericstore

import (
	"context"
	"salon_be/common"
)

func (s *genericStore) Create(ctx context.Context, modelName string, data interface{}) error {
	if err := s.db.Table(modelName).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
