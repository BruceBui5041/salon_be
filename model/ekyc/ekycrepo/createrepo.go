package ekycrepo

import (
	"context"
	"salon_be/common"
	"salon_be/component/ekycclient"
	models "salon_be/model"
	"salon_be/model/ekyc/ekycmodel"
)

type KYCStore interface {
	CreateKYCProfile(ctx context.Context, data *models.KYCProfile) error
	CreateDocument(ctx context.Context, data *models.IDDocument) error
	CreateFaceVerification(ctx context.Context, data *models.FaceVerification) error
}

type uploadRepo interface {
	UploadKYCImage(ctx context.Context, userId uint32, input *ekycmodel.UploadRequest) (*ekycmodel.KYCImageUploadRes, error)
}

type createKYCRepo struct {
	store      KYCStore
	ekycClient *ekycclient.EKYCClient
	uploadRepo uploadRepo
}

func NewCreateKYCRepo(store KYCStore, ekycClient *ekycclient.EKYCClient, uploadRepo uploadRepo) *createKYCRepo {
	return &createKYCRepo{
		store:      store,
		ekycClient: ekycClient,
		uploadRepo: uploadRepo,
	}
}

func (repo *createKYCRepo) ProcessKYCProfile(ctx context.Context, input *ekycmodel.CreateKYCProfileRequest) error {
	// Upload documents using uploadRepo
	frontImageRes, err := repo.uploadRepo.UploadKYCImage(ctx, input.UserID, &ekycmodel.UploadRequest{
		Image: input.FrontDocument,
	})
	if err != nil {
		return common.ErrInternal(err)
	}

	backImageRes, err := repo.uploadRepo.UploadKYCImage(ctx, input.UserID, &ekycmodel.UploadRequest{
		Image: input.BackDocument,
	})
	if err != nil {
		return common.ErrInternal(err)
	}

	faceImageRes, err := repo.uploadRepo.UploadKYCImage(ctx, input.UserID, &ekycmodel.UploadRequest{
		Image: input.FaceImage,
	})
	if err != nil {
		return common.ErrInternal(err)
	}

	// Document verification - Reduced to single classify call
	docType, err := repo.ekycClient.ClassifyDocument(frontImageRes.Object.Hash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Extract information and validate in one call
	docInfo, err := repo.ekycClient.ExtractDocumentInfo(frontImageRes.Object.Hash, backImageRes.Object.Hash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Face verification - Combined with liveness
	faceVerification, err := repo.ekycClient.VerifyFace(frontImageRes.Object.Hash, faceImageRes.Object.Hash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Create KYC profile
	profile := &models.KYCProfile{
		SQLModel:      common.SQLModel{Status: common.StatusActive},
		UserID:        input.UserID,
		CardID:        docInfo.Object.ID,
		FullName:      docInfo.Object.Name,
		DOB:           docInfo.Object.BirthDay,
		Gender:        docInfo.Object.Gender,
		Address:       docInfo.Object.RecentLocation,
		ClientSession: input.ClientSession,
	}

	if err := repo.store.CreateKYCProfile(ctx, profile); err != nil {
		return err
	}

	// Create document record
	document := &models.IDDocument{
		SQLModel:     common.SQLModel{Status: common.StatusActive},
		KYCProfileID: profile.Id,
		Type:         docType.Object.Type,
		Name:         docInfo.Object.Name,
		CardType:     docInfo.Object.CardType,
		ID:           docInfo.Object.ID,
		IDProbs:      docInfo.Object.IDProbs,
		BirthDay:     docInfo.Object.BirthDay,
		BirthDayProb: docInfo.Object.BirthDayProb,
	}

	if err := repo.store.CreateDocument(ctx, document); err != nil {
		return err
	}

	// Create face verification record
	faceData := &models.FaceVerification{
		SQLModel:     common.SQLModel{Status: common.StatusActive},
		KYCProfileID: profile.Id,
		Result:       faceVerification.Object.Result,
		Msg:          faceVerification.Object.Msg,
		Prob:         faceVerification.Object.Prob,
	}

	if err := repo.store.CreateFaceVerification(ctx, faceData); err != nil {
		return err
	}

	return nil
}
