package genericstore

import (
	"context"
	"fmt"
	"reflect"
	"salon_be/common"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/component/genericapi/modelhelper"
	"salon_be/component/modelsearch"

	"github.com/jinzhu/copier"
)

func (s *genericStore) Search(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
	result interface{},
) error {
	db := s.db

	modelType, err := modelhelper.GetModelType(input.Model)
	if err != nil {
		return common.ErrInternal(err)
	}

	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("result must be a pointer to a slice")
	}

	// Create a slice of the model type for querying
	sliceType := reflect.SliceOf(modelType)
	queryResult := reflect.New(sliceType).Interface()

	// Create a new instance of the model type for the query
	modelInstance := reflect.New(modelType).Interface()

	query := modelsearch.Search(
		ctx,
		db.Model(modelInstance),
		input.Model,
		input.Conditions,
		input.Fields,
		input.OrderBy,
	)

	// Perform the query
	err = query.Find(queryResult).Error
	if err != nil {
		return common.ErrDB(err)
	}

	copier.Copy(result, queryResult)

	return nil
}
