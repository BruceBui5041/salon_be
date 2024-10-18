package commentrepo

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/comment/commentmodel"

	"gorm.io/gorm"
)

type CreateCommentStore interface {
	CreateNewComment(
		ctx context.Context,
		newComment *models.Comment,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Comment, error)
}

type EnrollmentStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Enrollment, error)
}

type createCommentRepo struct {
	commentStore    CreateCommentStore
	enrollmentStore EnrollmentStore
}

func NewCreateCommentRepo(commentStore CreateCommentStore, enrollmentStore EnrollmentStore) *createCommentRepo {
	return &createCommentRepo{
		commentStore:    commentStore,
		enrollmentStore: enrollmentStore,
	}
}

func (repo *createCommentRepo) CreateNewComment(
	ctx context.Context,
	input *commentmodel.CreateComment,
) (*models.Comment, error) {
	courseUID, err := common.FromBase58(input.ServiceID)
	if err != nil {
		return nil, err
	}

	newComment := &models.Comment{
		UserID:    input.UserID,
		ServiceID: courseUID.GetLocalID(),
		Rate:      input.Rate,
		Content:   input.Content,
	}

	commentId, err := repo.commentStore.CreateNewComment(ctx, newComment)
	if err != nil {
		return nil, err
	}

	comment, err := repo.commentStore.FindOne(ctx, map[string]interface{}{"id": commentId})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (repo *createCommentRepo) ValidateCreateComment(ctx context.Context, userID, courseID uint32) (*models.Enrollment, error) {
	comment, err := repo.commentStore.FindOne(ctx, map[string]interface{}{
		"user_id":   userID,
		"course_id": courseID,
	})

	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return nil, common.ErrDB(err)
	}

	if comment != nil {
		return nil, common.ErrEntityExisted(models.CommentEntityName, errors.New("already commented"))
	}

	enrollment, err := repo.enrollmentStore.FindOne(ctx, map[string]interface{}{
		"user_id":   userID,
		"course_id": courseID,
	},
		"Payment",
	)

	if err != nil {
		return nil, common.ErrEntityNotFound(models.EnrollmentEntityName, err)
	}

	return enrollment, nil
}
