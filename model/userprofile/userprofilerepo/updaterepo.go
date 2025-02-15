package userprofilerepo

import (
	"context"
	"fmt"
	"salon_be/appconst"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/userprofile/userprofilemodel"
	"salon_be/storagehandler"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UpdateProfileStore interface {
	UpdateProfile(
		ctx context.Context,
		profileId uint32,
		data *models.UserProfile,
	) error
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.UserProfile, error)
}

type updateProfileRepo struct {
	store    UpdateProfileStore
	s3Client *s3.S3
}

func NewUpdateProfileRepo(store UpdateProfileStore, s3Client *s3.S3) *updateProfileRepo {
	return &updateProfileRepo{store: store, s3Client: s3Client}
}

func (repo *updateProfileRepo) UpdateProfile(
	ctx context.Context,
	userId string,
	profileId uint32,
	input *userprofilemodel.UpdateProfileModel,
) error {
	updatedUserProfile := &models.UserProfile{}
	if input.ProfilePictureURL != nil {
		pictureFile, err := input.ProfilePictureURL.Open()
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to open profile picture file", zap.Error(err))
			return fmt.Errorf("failed to open profile picture file: %w", err)
		}
		defer pictureFile.Close()

		key := storagehandler.GenerateUserProfilePictureS3Key(userId, input.ProfilePictureURL.Filename)

		err = storagehandler.UploadFileToS3(ctx, repo.s3Client, pictureFile, appconst.AWSPublicBucket, key)
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to upload profile picture to S3", zap.Error(err))
			return fmt.Errorf("failed to upload profile picture to S3: %w", err)
		}

		updatedUserProfile.ProfilePictureURL = key
	}

	if err := copier.Copy(&updatedUserProfile, input); err != nil {
		logger.AppLogger.Error(ctx, "Failed to copy updated user profile", zap.Error(err))
		return err
	}

	return repo.store.UpdateProfile(ctx, profileId, updatedUserProfile)
}

func (repo *updateProfileRepo) FindProfile(ctx context.Context, conditions map[string]interface{}) (*models.UserProfile, error) {
	return repo.store.FindOne(ctx, conditions)
}
