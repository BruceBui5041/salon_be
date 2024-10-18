package userprofilebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/userprofile/userprofilemodel"
	"salon_be/watermill"
	"salon_be/watermill/messagemodel"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type UpdateProfileRepo interface {
	UpdateProfile(ctx context.Context, userId string, profileId uint32, input *userprofilemodel.UpdateProfileModel) error
	FindProfile(ctx context.Context, conditions map[string]interface{}) (*models.UserProfile, error)
}

type updateProfileBiz struct {
	repo UpdateProfileRepo
}

func NewUpdateProfileBiz(repo UpdateProfileRepo) *updateProfileBiz {
	return &updateProfileBiz{repo: repo}
}

func (biz *updateProfileBiz) UpdateProfile(
	ctx context.Context,
	localPublisher *gochannel.GoChannel,
	input *userprofilemodel.UpdateProfileModel,
) error {
	requester, ok := ctx.Value(common.CurrentUser).(common.Requester)
	if !ok {
		return common.ErrNoPermission(errors.New("user not found"))
	}

	profile, err := biz.repo.FindProfile(ctx, map[string]interface{}{"user_id": requester.GetUserId()})
	if err != nil {
		return common.ErrCannotGetEntity(models.UserProfileEntityName, err)
	}

	if profile == nil {
		return common.ErrEntityNotFound(models.UserProfileEntityName, errors.New("profile not found"))
	}

	requester.Mask(false)
	if err := biz.repo.UpdateProfile(ctx, requester.GetFakeId(), profile.Id, input); err != nil {
		return common.ErrCannotUpdateEntity(models.UserProfileEntityName, err)
	}

	requester.Mask(false)
	if err := watermill.PublishUserUpdated(
		ctx,
		localPublisher,
		&messagemodel.UserUpdatedMessage{UserId: requester.GetFakeId()},
	); err != nil {
		return common.ErrInternal(err)
	}

	return nil
}
