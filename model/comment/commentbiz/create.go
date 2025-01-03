package commentbiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/comment/commentmodel"

	"github.com/jinzhu/copier"
)

type CommentRepo interface {
	CreateNewComment(ctx context.Context, input *commentmodel.CreateComment) (*models.Comment, error)
	ValidateCreateComment(ctx context.Context, userID, serviceID uint32) (*models.Enrollment, error)
}

type createCommentBiz struct {
	repo CommentRepo
}

func NewCreateCommentBiz(repo CommentRepo) *createCommentBiz {
	return &createCommentBiz{repo: repo}
}

func (c *createCommentBiz) CreateNewComment(ctx context.Context, input *commentmodel.CreateComment) error {
	if input.UserID == 0 {
		return common.ErrInvalidRequest(errors.New("user ID is required"))
	}

	if input.ServiceID == "" {
		return common.ErrInvalidRequest(errors.New("service ID is required"))
	}

	if input.Rate > 5 || input.Rate < 1 {
		return common.ErrInvalidRequest(
			errors.New("rate must be equal and greater than 1 and equal and lesser than 5"),
		)
	}

	if input.Content == "" {
		return common.ErrInvalidRequest(errors.New("comment content is required"))
	}

	if len(input.Content) > 1000 {
		return common.ErrInvalidRequest(errors.New("comment content must not exceed 1000 characters"))
	}

	serviceUID, err := common.FromBase58(input.ServiceID)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Check if the user is enrolled in the service and the payment is completed
	enrollment, err := c.repo.ValidateCreateComment(ctx, input.UserID, serviceUID.GetLocalID())
	if err != nil {
		return err
	}

	if enrollment == nil || enrollment.Payment == nil || enrollment.Payment.TransactionStatus != "completed" {
		return common.ErrNoPermission(errors.New("user is not enrolled or payment is not completed"))
	}

	comment, err := c.repo.CreateNewComment(ctx, input)
	if err != nil {
		return common.ErrCannotCreateEntity(models.CommentEntityName, err)
	}

	var res commentmodel.CommentRes
	copier.Copy(&res, comment)

	return nil
}
