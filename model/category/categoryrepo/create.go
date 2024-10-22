package categoryrepo

import (
	"context"
	"fmt"
	"salon_be/appconst"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/category/categorymodel"
	"salon_be/storagehandler"

	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type CreateCategoryStore interface {
	Create(
		ctx context.Context,
		newCategory *models.Category,
	) (*models.Category, error)
	Update(
		ctx context.Context,
		id uint32,
		data *models.Category,
	) error

	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
}

type createCategoryRepo struct {
	store    CreateCategoryStore
	s3Client *s3.S3
}

func NewCreateCategoryRepo(store CreateCategoryStore, s3Client *s3.S3) *createCategoryRepo {
	return &createCategoryRepo{
		store:    store,
		s3Client: s3Client,
	}
}

func (repo *createCategoryRepo) CreateNewCategory(
	ctx context.Context,
	input *categorymodel.CreateCategory,
) (*models.Category, error) {
	newCategory := &models.Category{
		Name:        input.Name,
		Description: input.Description,
		Code:        input.Code,
	}

	if input.ParentID != nil && *input.ParentID != "" {
		uid, err := common.FromBase58(*input.ParentID)
		if err != nil {
			return nil, common.ErrInvalidRequest(err)
		}

		localID := uid.GetLocalID()
		newCategory.ParentID = &localID
	}

	category, err := repo.store.Create(ctx, newCategory)
	if err != nil {
		return nil, common.ErrDB(err)
	}

	if input.Image != nil {
		pictureFile, err := input.Image.Open()
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to open image file", zap.Error(err))
			return nil, fmt.Errorf("failed to open image file: %w", err)
		}
		defer pictureFile.Close()

		key := storagehandler.GenerateCategoryImageS3Key(category.GetFakeId(), input.Image.Filename)

		err = storagehandler.UploadFileToS3(ctx, repo.s3Client, pictureFile, appconst.AWSPublicBucket, key)
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to upload image to S3", zap.Error(err))
			return nil, fmt.Errorf("failed to upload image to S3: %w", err)
		}

		err = repo.store.Update(ctx, category.Id, &models.Category{Image: key})
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to update category", zap.Error(err))
			return nil, fmt.Errorf("failed to update category: %w", err)
		}
	}

	return category, nil
}

func (repo *createCategoryRepo) GetCategory(ctx context.Context, id uint32) (*models.Category, error) {
	category, err := repo.store.FindOne(ctx, map[string]interface{}{
		"id": id,
	}, "Parent")

	if err != nil {
		return nil, common.ErrDB(err)
	}

	return category, nil
}
