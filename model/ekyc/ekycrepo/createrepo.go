package ekycrepo

import (
	"context"
	"salon_be/common"
	"salon_be/component/ekycclient"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/ekyc/ekycmodel"

	"go.uber.org/zap"
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
	logger.AppLogger.Info(ctx, "starting KYC profile processing",
		zap.Uint32("user_id", input.UserID))

	// Upload documents using uploadRepo
	frontImageRes, err := repo.uploadRepo.UploadKYCImage(ctx, input.UserID, &ekycmodel.UploadRequest{
		Image: input.FrontDocument,
	})
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to upload front image",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
		return common.ErrInternal(err)
	}

	backImageRes, err := repo.uploadRepo.UploadKYCImage(ctx, input.UserID, &ekycmodel.UploadRequest{
		Image: input.BackDocument,
	})
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to upload back image",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
		return common.ErrInternal(err)
	}

	faceImageRes, err := repo.uploadRepo.UploadKYCImage(ctx, input.UserID, &ekycmodel.UploadRequest{
		Image: input.FaceImage,
	})
	if err != nil {
		logger.AppLogger.Error(ctx, "failed to upload face image",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
		return common.ErrInternal(err)
	}

	logger.AppLogger.Info(ctx, "document verification starting",
		zap.Uint32("user_id", input.UserID))

	// Document verification
	docType, err := repo.ekycClient.ClassifyDocument(ctx, frontImageRes.Object.Hash, input.ClientSession)
	if err != nil {
		logger.AppLogger.Error(ctx, "document classification failed",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
		return common.ErrInternal(err)
	}

	// Extract information and validate
	docInfo, err := repo.ekycClient.ExtractDocumentInfo(ctx, frontImageRes.Object.Hash, backImageRes.Object.Hash, input.ClientSession)
	if err != nil {
		logger.AppLogger.Error(ctx, "document info extraction failed",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
		return common.ErrInternal(err)
	}

	logger.AppLogger.Info(ctx, "face verification starting",
		zap.Uint32("user_id", input.UserID))

	// Face verification
	faceVerification, err := repo.ekycClient.VerifyFace(ctx, frontImageRes.Object.Hash, faceImageRes.Object.Hash, input.ClientSession)
	if err != nil {
		logger.AppLogger.Error(ctx, "face verification failed",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
		return common.ErrInternal(err)
	}

	logger.AppLogger.Info(ctx, "creating KYC profile records",
		zap.Uint32("user_id", input.UserID))

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
		logger.AppLogger.Error(ctx, "failed to create KYC profile",
			zap.Uint32("user_id", input.UserID),
			zap.Error(err))
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
		Nationality:  docInfo.Object.Nationality,
		Gender:       docInfo.Object.Gender,
		ValidDate:    docInfo.Object.ValidDate,
		IssueDate:    docInfo.Object.IssueDate,
		IssuePlace:   docInfo.Object.IssuePlace,
	}

	if err := repo.store.CreateDocument(ctx, document); err != nil {
		logger.AppLogger.Error(ctx, "failed to create document record",
			zap.Uint32("user_id", input.UserID),
			zap.String("doc_id", document.ID),
			zap.Error(err))
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
		logger.AppLogger.Error(ctx, "failed to create face verification record",
			zap.Uint32("user_id", input.UserID),
			zap.Uint32("profile_id", profile.Id),
			zap.Error(err))
		return err
	}

	logger.AppLogger.Info(ctx, "KYC profile processing completed successfully",
		zap.Uint32("user_id", input.UserID),
		zap.Uint32("profile_id", profile.Id))

	return nil
}
