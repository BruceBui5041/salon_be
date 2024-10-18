package commentbiz

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/comment/commentmodel"
)

type UpdateCommentRepo interface {
	UpdateComment(ctx context.Context, id uint32, input *commentmodel.UpdateComment) error
}

type updateCommentBiz struct {
	repo UpdateCommentRepo
}

func NewUpdateCommentBiz(repo UpdateCommentRepo) *updateCommentBiz {
	return &updateCommentBiz{repo: repo}
}

func (c *updateCommentBiz) UpdateComment(ctx context.Context, id uint32, input *commentmodel.UpdateComment) error {
	if input.Content != "" && len(input.Content) > 1000 {
		return common.ErrInvalidRequest(common.NewCustomError(nil, "comment content must not exceed 1000 characters", "ErrCommentContentTooLong"))
	}

	if err := c.repo.UpdateComment(ctx, id, input); err != nil {
		return common.ErrCannotUpdateEntity(models.CommentEntityName, err)
	}

	return nil
}
