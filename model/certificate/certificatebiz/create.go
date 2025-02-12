package certificatebiz

import (
	"context"
	"fmt"
	"salon_be/common"
	"salon_be/component/logger"
	models "salon_be/model"
	"salon_be/model/certificate/certificateerror"
	"salon_be/model/certificate/certificatemodel"
	"strings"
	"time"

	"go.uber.org/zap"
)

type CreateRepository interface {
	Create(ctx context.Context, data *models.Certificate) error
}

type createBiz struct {
	repo CreateRepository
}

func NewCreateBiz(repo CreateRepository) *createBiz {
	return &createBiz{repo: repo}
}

func (biz *createBiz) CreateCertificate(ctx context.Context, data *certificatemodel.CreateCertificateInput) error {
	// Get file info from uploaded file
	fileHeader := data.File

	// Validate file type using the original filename
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".pdf") {
		return certificateerror.ErrInvalidFileType(fmt.Errorf("invalid file type: %s", fileHeader.Filename))
	}

	// Validate file size (5MB = 5 * 1024 * 1024 bytes)
	if fileHeader.Size > 5*1024*1024 {
		return certificateerror.ErrFileTooLarge(fmt.Errorf("file size exceeds 5MB limit"))
	}

	tempUser := common.SQLModel{Id: data.OwnerID}
	tempUser.GenUID(common.DbTypeUser)

	s3Key := GenerateCertificateS3KeyWithTimestamp(tempUser.GetFakeId(), fileHeader.Filename)

	certificate := &models.Certificate{
		URL:       s3Key,
		Type:      data.Type,
		OwnerID:   data.OwnerID,
		CreatorID: data.CreatorID,
	}

	if err := biz.repo.Create(ctx, certificate); err != nil {
		logger.AppLogger.Error(ctx, "create certificate failed", zap.Error(err))
		return common.ErrCannotCreateEntity(models.CertificateEntityName, err)
	}

	return nil
}

func GenerateCertificateS3KeyWithTimestamp(ownerID string, filename string) string {
	timestamp := time.Now().UTC().Unix()
	return fmt.Sprintf(
		"certificates/%s/%d_%s",
		ownerID,
		timestamp,
		filename,
	)
}
