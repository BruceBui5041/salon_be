package userbiz

import (
	"context"
	"video_server/common"
	"video_server/component/genericapi/genericmodel"
	models "video_server/model"
	"video_server/model/user/usermodel"

	"github.com/jinzhu/copier"
)

type SearchUserStore interface {
	Search(
		ctx context.Context,
		input genericmodel.SearchModelRequest,
	) ([]*models.User, error)
}

type searchUserBiz struct {
	store SearchUserStore
}

func NewSearchUserBiz(store SearchUserStore) *searchUserBiz {
	return &searchUserBiz{store: store}
}

func (b *searchUserBiz) SearchUsers(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
) ([]*usermodel.GetUserResponse, error) {
	users, err := b.store.Search(ctx, input)
	if err != nil {
		panic(common.ErrDB(err))
	}

	var result []*usermodel.GetUserResponse
	if err := copier.Copy(&result, users); err != nil {
		panic(common.ErrInternal(err))
	}

	for _, user := range result {
		user.Mask(false)
	}

	return result, nil
}
