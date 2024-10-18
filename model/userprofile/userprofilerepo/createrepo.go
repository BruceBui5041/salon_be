package userprofilerepo

import (
	"context"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/userprofile/userprofilemodel"
)

type CreateUserProfileStore interface {
	CreateNewUserProfile(
		ctx context.Context,
		newUserProfile *models.UserProfile,
	) (uint32, error)
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.UserProfile, error)
}

type createUserProfileRepo struct {
	userProfileStore CreateUserProfileStore
}

func NewCreateUserProfileRepo(userProfileStore CreateUserProfileStore) *createUserProfileRepo {
	return &createUserProfileRepo{
		userProfileStore: userProfileStore,
	}
}

func (repo *createUserProfileRepo) CreateNewUserProfile(
	ctx context.Context,
	input *userprofilemodel.CreateUserProfile,
) (*models.UserProfile, error) {
	newUserProfile := &models.UserProfile{
		UserID:      input.UserID,
		PhoneNumber: input.PhoneNumber,
		Occupation:  input.Occupation,
		Biography:   input.Biography,
		LinkedIn:    input.LinkedIn,
		Facebook:    input.Facebook,
		Twitter:     input.Twitter,
		Instagram:   input.Instagram,
	}

	userProfileId, err := repo.userProfileStore.CreateNewUserProfile(ctx, newUserProfile)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.UserProfileEntityName, err)
	}

	userProfile, err := repo.userProfileStore.FindOne(ctx, map[string]interface{}{"id": userProfileId})
	if err != nil {
		return nil, common.ErrCannotGetEntity(models.UserProfileEntityName, err)
	}

	return userProfile, nil
}
