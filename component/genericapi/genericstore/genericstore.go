package genericstore

import (
	"context"
	"reflect"
	"salon_be/common"
	"salon_be/component/genericapi/genericmodel"

	"gorm.io/gorm"
)

type GenericStoreInterface interface {
	Search(ctx context.Context, input genericmodel.SearchModelRequest, result interface{}) error
	Create(ctx context.Context, modelName string, data interface{}) error
}

type genericStore struct {
	db *gorm.DB
}

func NewGenericStore(db *gorm.DB) *genericStore {
	return &genericStore{db: db}
}

func (s *genericStore) FindOne(ctx context.Context, model interface{}, conditions map[string]interface{}, moreInfo ...string) error {
	db := s.db

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if err := db.Where(conditions).First(model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return common.RecordNotFound
		}
		return common.ErrDB(err)
	}

	return nil
}

func (s *genericStore) Update(ctx context.Context, model interface{}) error {
	if err := s.db.Save(model).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *genericStore) Delete(ctx context.Context, model interface{}) error {
	if err := s.db.Delete(model).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (s *genericStore) List(ctx context.Context, model interface{}, filter map[string]interface{}, paging *common.Paging, moreInfo ...string) error {
	db := s.db

	db = db.Table(reflect.TypeOf(model).Elem().Name())

	if f := filter; len(f) > 0 {
		for k, v := range f {
			db = db.Where(k, v)
		}
	}

	if err := db.Count(&paging.Total).Error; err != nil {
		return common.ErrDB(err)
	}

	for i := range moreInfo {
		db = db.Preload(moreInfo[i])
	}

	if v := paging.FakeCursor; v != "" {
		if uid, err := common.FromBase58(v); err == nil {
			db = db.Where("id < ?", uid.GetLocalID())
		}
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db = db.Offset(int(offset))
	}

	if err := db.
		Limit(int(paging.Limit)).
		Order("id desc").
		Find(model).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}

func (s *genericStore) Count(ctx context.Context, model interface{}, conditions map[string]interface{}) (int64, error) {
	var count int64
	if err := s.db.Model(model).Where(conditions).Count(&count).Error; err != nil {
		return 0, common.ErrDB(err)
	}
	return count, nil
}
