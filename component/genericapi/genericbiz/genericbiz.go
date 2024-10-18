package genericbiz

import (
	"context"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/component/genericapi/genericstore"
)

type GenericBiz interface {
	Create(ctx context.Context, input genericmodel.CreateRequest) (interface{}, error)
	Search(
		ctx context.Context,
		input genericmodel.SearchModelRequest,
		result interface{},
	) error
}

type genericBiz struct {
	store genericstore.GenericStoreInterface
}

func NewGenericBiz(store genericstore.GenericStoreInterface) *genericBiz {
	return &genericBiz{store: store}
}
