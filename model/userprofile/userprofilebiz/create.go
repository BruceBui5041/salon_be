package userprofilebiz

import (
	"context"
	"errors"
	"salon_be/common"
	models "salon_be/model"
	"salon_be/model/userprofile/userprofilemodel"
)

type UserProfileRepo interface {
	CreateNewUserProfile(ctx context.Context, input *userprofilemodel.CreateUserProfile) (*models.UserProfile, error)
}

type createUserProfileBiz struct {
	repo UserProfileRepo
}

func NewCreateUserProfileBiz(repo UserProfileRepo) *createUserProfileBiz {
	return &createUserProfileBiz{repo: repo}
}

func (c *createUserProfileBiz) CreateNewUserProfile(ctx context.Context, input *userprofilemodel.CreateUserProfile) (*models.UserProfile, error) {
	if input.UserID == 0 {
		return nil, errors.New("user ID is required")
	}
	if len(input.FirstName) > 50 {
		return nil, errors.New("first name must not exceed 50 characters")
	}
	if len(input.LastName) > 50 {
		return nil, errors.New("last name must not exceed 50 characters")
	}
	if len(input.PhoneNumber) > 20 {
		return nil, errors.New("phone number must not exceed 20 characters")
	}
	if len(input.Occupation) > 100 {
		return nil, errors.New("occupation must not exceed 100 characters")
	}
	if len(input.LinkedIn) > 255 {
		return nil, errors.New("LinkedIn URL must not exceed 255 characters")
	}
	if len(input.Facebook) > 255 {
		return nil, errors.New("facebook url must not exceed 255 characters")
	}
	if len(input.Twitter) > 255 {
		return nil, errors.New("twitter url must not exceed 255 characters")
	}
	if len(input.Instagram) > 255 {
		return nil, errors.New("instagram url must not exceed 255 characters")
	}

	userProfile, err := c.repo.CreateNewUserProfile(ctx, input)
	if err != nil {
		return nil, common.ErrCannotCreateEntity(models.UserProfileEntityName, err)
	}
	return userProfile, nil
}
