package genericbiz

import (
	"context"
	"salon_be/common"
	"salon_be/component/genericapi/genericmodel"
	"salon_be/component/genericapi/modelhelper"
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
