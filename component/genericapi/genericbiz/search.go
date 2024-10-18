package genericbiz

import (
	"context"
	"salon_be/common"
	"salon_be/component/genericapi/genericmodel"
)

func (b *genericBiz) Search(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
	result interface{},
) error {
	if err := b.store.Search(ctx, input, result); err != nil {
		return common.ErrDB(err)
	}

	return nil
}
