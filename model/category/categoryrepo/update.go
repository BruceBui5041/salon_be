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
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type UpdateCategoryStore interface {
	FindOne(
		ctx context.Context,
		conditions map[string]interface{},
		moreInfo ...string,
	) (*models.Category, error)
	Update(
		ctx context.Context,
		id uint32,
		data *models.Category,
	) error

	UpdateParentId(
		ctx context.Context,
		id uint32,
		parentID *uint32,
	) error
}

type updateCategoryRepo struct {
	store    UpdateCategoryStore
	s3Client *s3.S3
}

func NewUpdateCategoryRepo(
	store UpdateCategoryStore,
	s3Client *s3.S3,
) *updateCategoryRepo {
	return &updateCategoryRepo{
		store:    store,
		s3Client: s3Client,
	}
}

func (repo *updateCategoryRepo) GetCategory(ctx context.Context, id uint32) (*models.Category, error) {
	category, err := repo.store.FindOne(ctx, map[string]interface{}{"id": id})
	if err != nil {
		return nil, common.ErrDB(err)
	}

	return category, nil
}

func (repo *updateCategoryRepo) UpdateCategory(ctx context.Context, id uint32, data *categorymodel.UpdateCategory) error {
	var categ models.Category
	copier.Copy(&categ, data)

	if data.ParentID != nil {
		if *data.ParentID == "" {
			categ.ParentID = nil
			if err := repo.store.UpdateParentId(ctx, id, nil); err != nil {
				return common.ErrDB(err)
			}
		} else {
			uid, err := common.FromBase58(*data.ParentID)
			if err != nil {
				return common.ErrInvalidRequest(err)
			}

			localID := uid.GetLocalID()
			categ.ParentID = &localID
		}
	}

	if data.Image != nil {
		pictureFile, err := data.Image.Open()
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to open image file", zap.Error(err))
			return fmt.Errorf("failed to open image file: %w", err)
		}
		defer pictureFile.Close()

		sqlModel := common.SQLModel{Id: id}
		sqlModel.GenUID(common.DBTypeCategory)

		key := storagehandler.GenerateCategoryImageS3Key(sqlModel.GetFakeId(), data.Image.Filename)

		err = storagehandler.UploadFileToS3(ctx, repo.s3Client, pictureFile, appconst.AWSPublicBucket, key)
		if err != nil {
			logger.AppLogger.Error(ctx, "Failed to upload image to S3", zap.Error(err))
			return fmt.Errorf("failed to upload image to S3: %w", err)
		}

		categ.Image = key

	}

	oldCateg, err := repo.store.FindOne(ctx, map[string]interface{}{
		"id": id,
	})

	if err != nil {
		logger.AppLogger.Error(ctx, "Failed to find old category", zap.Error(err))
		return err
	}

	if err := repo.store.Update(ctx, id, &categ); err != nil {
		return common.ErrDB(err)
	}

	if categ.Image != "" && oldCateg.Image != "" {
		if err := storagehandler.RemoveFileFromS3(
			ctx,
			repo.s3Client,
			appconst.AWSPublicBucket,
			oldCateg.OriginImage,
		); err != nil {
			logger.AppLogger.Error(ctx, "Failed to remove old profile picture from S3", zap.Error(err))
		}
	}

	return nil
}
