package genericbiz

import (
	"context"
	"video_server/common"
	"video_server/component/genericapi/genericmodel"
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
