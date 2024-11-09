package otpstore

import (
	"context"
	models "salon_be/model"
)

func (s *sqlStore) FindOne(
	ctx context.Context,
	conditions map[string]interface{},
	moreKeys ...string,
) (*models.OTP, error) {
	var result models.OTP
	db := s.db

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	query := db
	if expiresAtCond, exists := conditions["expires_at > ?"]; exists {
		query = query.Where("expires_at > ?", expiresAtCond)
		delete(conditions, "expires_at > ?")
	}

	if err := query.Where(conditions).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil

}
