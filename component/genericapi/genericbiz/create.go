package genericbiz

import (
	"context"
	"video_server/common"
	"video_server/component/genericapi/genericmodel"
	"video_server/component/genericapi/modelhelper"
)

func (biz *genericBiz) Create(ctx context.Context, input genericmodel.CreateRequest) (interface{}, error) {
	if _, err := modelhelper.GetModelType(input.Model); err != nil {
		return nil, common.ErrInvalidRequest(err)
	}

	if err := biz.store.Create(ctx, input.Model, input.Data); err != nil {
		return nil, common.ErrCannotCreateEntity(input.Model, err)
	}

	return input.Data, nil
}
