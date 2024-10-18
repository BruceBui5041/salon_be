package genericbiz

import (
	"context"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/genericapi/genericstore"
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
