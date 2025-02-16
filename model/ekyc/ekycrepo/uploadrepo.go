package ekycrepo

import (
	"context"
	"mime/multipart"
	"salon_be/common"
	"salon_be/component/ekycclient"
	models "salon_be/model"
	"salon_be/model/ekyc/ekycmodel"
	"salon_be/storagehandler"
)

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

type KYCImageStore interface {
	CreateKYCImage(ctx context.Context, data *models.KYCImageUpload) error
	UpdateKYCImage(ctx context.Context, kycImageId uint32, data *models.KYCImageUpload) error
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
	// Upload to VNPT eKYC service
	ekycRes, err := repo.ekycClient.UploadFile(ctx, input.Image)
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
		Provider:     "vnpt",
	}

	if err := repo.store.CreateKYCImage(ctx, kycImage); err != nil {
		return nil, common.ErrDB(err)
	}

	tempObj := common.SQLModel{Id: kycImage.Id}
	tempObj.GenUID(common.DBTypeKYCImage)

	s3ObjectKey := storagehandler.GenerateKYCImageS3Key(tempObj.GetFakeId(), input.Image.Filename)

	image, err := repo.imageRepo.CreateImage(ctx, input.Image, 0, userId, s3ObjectKey, "kyc_image")
	if err != nil {
		return nil, err
	}

	if err := repo.store.UpdateKYCImage(
		ctx,
		kycImage.Id,
		&models.KYCImageUpload{ImageID: &image.Id},
	); err != nil {
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
