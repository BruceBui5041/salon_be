package commentbiz

import (
	"context"
	"salon_be/common"
	"salon_be/component/genericapi/genericmodel"
	models "salon_be/model"
	"salon_be/model/comment/commentmodel"

	"github.com/jinzhu/copier"
)

type SearchCommentStore interface {
	Search(
		ctx context.Context,
		input genericmodel.SearchModelRequest,
	) ([]*models.Comment, error)
}

type searchCommentBiz struct {
	store SearchCommentStore
}

func NewSearchCommentBiz(store SearchCommentStore) *searchCommentBiz {
	return &searchCommentBiz{store: store}
}

func (b *searchCommentBiz) SearchComments(
	ctx context.Context,
	input genericmodel.SearchModelRequest,
) ([]*commentmodel.CommentRes, error) {
	comments, err := b.store.Search(ctx, input)
	if err != nil {
		return nil, common.ErrDB(err)
	}
	var res []*commentmodel.CommentRes
	copier.Copy(&res, comments)

	return res, nil
}
