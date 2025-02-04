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
type createKYCRepo struct {
	store      KYCStore
	ekycClient *ekycclient.EKYCClient
}

func NewCreateKYCRepo(store KYCStore, ekycClient *ekycclient.EKYCClient) *createKYCRepo {
	return &createKYCRepo{
		store:      store,
		ekycClient: ekycClient,
	}
}

func (repo *createKYCRepo) ProcessKYCProfile(ctx context.Context, input *ekycmodel.CreateKYCProfileRequest) error {
	// Upload documents
	frontHash, err := repo.ekycClient.UploadFile(input.FrontDocument)
	if err != nil {
		return common.ErrInternal(err)
	}

	backHash, err := repo.ekycClient.UploadFile(input.BackDocument)
	if err != nil {
		return common.ErrInternal(err)
	}

	faceHash, err := repo.ekycClient.UploadFile(input.FaceImage)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Document verification
	docType, err := repo.ekycClient.ClassifyDocument(frontHash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	docLiveness, err := repo.ekycClient.ValidateDocument(frontHash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Extract information
	docInfo, err := repo.ekycClient.ExtractDocumentInfo(frontHash, backHash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Face verification
	faceVerification, err := repo.ekycClient.VerifyFace(frontHash, faceHash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	faceLiveness, err := repo.ekycClient.CheckFaceLiveness(faceHash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	maskCheck, err := repo.ekycClient.CheckFaceMask(faceHash, input.ClientSession)
	if err != nil {
		return common.ErrInternal(err)
	}

	// Create KYC profile
	profile := &models.KYCProfile{
		SQLModel:      common.SQLModel{Status: common.StatusInactive},
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
		SQLModel:       common.SQLModel{Status: common.StatusActive},
		KYCProfileID:   profile.Id,
		Type:           docType.Object.Type,
		Name:           docInfo.Object.Name,
		CardType:       docInfo.Object.CardType,
		ID:             docInfo.Object.ID,
		IDProbs:        docInfo.Object.IDProbs,
		BirthDay:       docInfo.Object.BirthDay,
		BirthDayProb:   docInfo.Object.BirthDayProb,
		LivenessStatus: docLiveness.Object.LivenessStatus,
		LivenessMsg:    docLiveness.Object.LivenessMsg,
		FaceSwapping:   docLiveness.Object.FaceSwapping,
		FakeLiveness:   docLiveness.Object.FakeLiveness,
	}

	if err := repo.store.CreateDocument(ctx, document); err != nil {
		return err
	}

	// Create face verification record
	faceData := &models.FaceVerification{
		SQLModel:       common.SQLModel{Status: common.StatusActive},
		KYCProfileID:   profile.Id,
		Result:         faceVerification.Object.Result,
		Msg:            faceVerification.Object.Msg,
		Prob:           faceVerification.Object.Prob,
		LivenessStatus: faceLiveness.Object.LivenessStatus,
		LivenessMsg:    faceLiveness.Object.LivenessMsg,
		IsEyeOpen:      faceLiveness.Object.IsEyeOpen,
		Masked:         maskCheck.Object.Masked,
	}

	if err := repo.store.CreateFaceVerification(ctx, faceData); err != nil {
		return err
	}

	return nil
}
