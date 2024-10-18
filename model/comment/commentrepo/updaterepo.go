package commentrepo

import (
	"context"
	"video_server/common"
	models "video_server/model"
	"video_server/model/comment/commentmodel"
)

type UpdateCommentStore interface {
	Update(
		ctx context.Context,
		id uint32,
		updateData *models.Comment,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Comment, error)
}

type updateCommentRepo struct {
	updateCommentStore UpdateCommentStore
}

func NewUpdateCommentRepo(
	updateCommentStore UpdateCommentStore,
) *updateCommentRepo {
	return &updateCommentRepo{
		updateCommentStore: updateCommentStore,
	}
}

func (repo *updateCommentRepo) UpdateComment(ctx context.Context, commentId uint32, input *commentmodel.UpdateComment) error {
	comment, err := repo.updateCommentStore.FindOne(ctx, map[string]interface{}{"id": commentId}, "User")
	if err != nil {
		return common.ErrCannotGetEntity(models.CommentEntityName, err)
	}

	requester := ctx.Value(common.CurrentUser).(common.Requester)

	if comment.UserID != requester.GetUserId() {
		return common.ErrNoPermission(common.NewCustomError(nil, "only the author can update the comment", "ErrCommentUpdateNoPermission"))
	}

	updateData := &models.Comment{
		Content: input.Content,
		Rate:    input.Rate,
	}

	if err := repo.updateCommentStore.Update(ctx, commentId, updateData); err != nil {
		return common.ErrCannotUpdateEntity(models.CommentEntityName, err)
	}

	return nil
}
