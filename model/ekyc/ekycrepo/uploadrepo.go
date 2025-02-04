package ekycrepo

import (
	"context"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/ekycclient"
	models "salon_be/model"
	"salon_be/model/ekyc/ekycmodel"
)

type ImageRepo interface {
	CreateImage(
		ctx context.Context,
		file *multipart.FileHeader,
		serviceID uint32,
		userID uint32,
	) (*models.Image, error)
}

type KYCImageStore interface {
	CreateKYCImage(ctx context.Context, data *models.KYCImageUpload) error
}

type kycImageUploadRepo struct {
	store      KYCImageStore
	ekycClient *ekycclient.EKYCClient
	imageRepo  ImageRepo
}

func NewKYCImageUploadRepo(
	store KYCImageStore,
	ekycClient *ekycclient.EKYCClient,
	imageRepo ImageRepo,
) *kycImageUploadRepo {
	return &kycImageUploadRepo{
		store:      store,
		ekycClient: ekycClient,
		imageRepo:  imageRepo,
	}
}

func (repo *kycImageUploadRepo) UploadKYCImage(
	ctx context.Context,
	userId uint32,
	input *ekycmodel.UploadRequest,
) (*ekycmodel.KYCImageUploadRes, error) {
	// Use image repo interface to handle image upload
	image, err := repo.imageRepo.CreateImage(ctx, input.Image, 0, userId)
	if err != nil {
		return nil, err
	}

	// Upload to VNPT eKYC service
	ekycRes, err := repo.ekycClient.UploadFile(input.Image)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	// Create KYC image record
	kycImage := &models.KYCImageUpload{
		SQLModel: common.SQLModel{
			Status: common.StatusActive,
		},
		UserID:       userId,
		FileName:     ekycRes.Object.FileName,
		Title:        ekycRes.Object.Title,
		Description:  ekycRes.Object.Description,
		Hash:         ekycRes.Object.Hash,
		FileType:     ekycRes.Object.FileType,
		UploadedDate: ekycRes.Object.UploadedDate,
		StorageType:  ekycRes.Object.StorageType,
		TokenId:      ekycRes.Object.TokenId,
		ImageID:      &image.Id,
		Provider:     "vnpt",
	}

	if err := repo.store.CreateKYCImage(ctx, kycImage); err != nil {
		return nil, common.ErrDB(err)
	}

	return &ekycmodel.KYCImageUploadRes{
		Message: ekycRes.Message,
		Object: ekycmodel.ImageUploadObject{
			FileName:     ekycRes.Object.FileName,
			Title:        ekycRes.Object.Title,
			Description:  ekycRes.Object.Description,
			Hash:         ekycRes.Object.Hash,
			FileType:     ekycRes.Object.FileType,
			UploadedDate: ekycRes.Object.UploadedDate,
			StorageType:  ekycRes.Object.StorageType,
			TokenId:      ekycRes.Object.TokenId,
		},
	}, nil
}
