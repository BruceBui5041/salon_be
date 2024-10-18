package coursestore

import (
	"video_server/component/genericapi/genericstore"

	"gorm.io/gorm"
)

type sqlStore struct {
	db *gorm.DB
}

func NewSQLStore(db *gorm.DB) *sqlStore {
	return &sqlStore{db: db}
}

func NewGenericStore(store genericstore.GenericStoreInterface) genericstore.GenericStoreInterface {
	return store
}
