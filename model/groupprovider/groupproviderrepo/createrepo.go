package groupproviderrepo

import (
	"context"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/storagehandler"

	"go.uber.org/zap"
)

type CreateGroupProviderStore interface {
	Create(ctx context.Context, data *models.GroupProvider) error
	FindOne(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*models.GroupProvider, error)
}

type CreateUserStore interface {
	FindOne(ctx context.Context, conditions map[string]interface{}, moreKeys ...string) (*models.User, error)
}

type CreateServiceStore interface {
	UpdateMany(ctx context.Context, condition map[string]interface{}, data *models.Service) error
}

type ImageRepo interface {
	CreateImage(
		ctx context.Context,
		file *multipart.FileHeader,
		groupProviderID uint32,
		userID uint32,
		s3ObjectKey string,
		refType string,
	) (*models.Image, error)
}

type createRepo struct {
	store        CreateGroupProviderStore
	userStore    CreateUserStore
	serviceStore CreateServiceStore
	imageRepo    ImageRepo
}

func NewCreateRepo(
	store CreateGroupProviderStore,
	userStore CreateUserStore,
	serviceStore CreateServiceStore,
	imageRepo ImageRepo,
) *createRepo {
	return &createRepo{
		store:        store,
		userStore:    userStore,
		serviceStore: serviceStore,
		imageRepo:    imageRepo,
	}
}

func (r *createRepo) FindUserByID(ctx context.Context, userID uint32) (*models.User, error) {
	user, err := r.userStore.FindOne(ctx, map[string]interface{}{"id": userID}, "Roles")
	if err != nil {
		return nil, common.ErrEntityNotFound(models.UserEntityName, err)
	}
	return user, nil
}

func (r *createRepo) FindGroupProviderByOwner(ctx context.Context, ownerID uint32) (*models.GroupProvider, error) {
	groupProvider, err := r.store.FindOne(ctx, map[string]interface{}{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	return groupProvider, nil
}

func (r *createRepo) Create(
	ctx context.Context,
	groupProvider *models.GroupProvider,
	images []*multipart.FileHeader,
	creatorID uint32,
) error {
	if err := r.store.Create(ctx, groupProvider); err != nil {
		return err
	}

	if len(images) > 0 {
		groupProvider.Mask(true)

		for _, file := range images {
			objectKey := storagehandler.GenerateGroupProviderImageS3Key(
				groupProvider.GetFakeId(), file.Filename,
			)

			img, err := r.imageRepo.CreateImage(
				ctx,
				file,
				groupProvider.Id,
				creatorID,
				objectKey,
				"group_provider",
			)
			if err != nil {
				logger.AppLogger.Error(ctx, "failed to upload group provider image", zap.Error(err))
				return err
			}
			groupProvider.Images = append(groupProvider.Images, img)
		}
	}

	return nil
}

func (r *createRepo) UpdateServicesGroupProvider(ctx context.Context, ownerID uint32, groupProviderID uint32) error {
	condition := map[string]interface{}{"owner_id": ownerID, "group_provider_id": nil}
	service := &models.Service{
		GroupProviderID: &groupProviderID,
	}

	if err := r.serviceStore.UpdateMany(ctx, condition, service); err != nil {
		return err
	}
	return nil
}
