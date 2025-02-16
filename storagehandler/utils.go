package storagehandler

import (
	"fmt"
	"path/filepath"
	"salon_be/appconst"
	"salon_be/utils"

	"github.com/google/uuid"
)

type VideoInfo struct {
	UploadedBy        string `json:"uploaded_by"`
	ServiceId         string `json:"service_id"`
	VideoId           string `json:"video_id"`
	ThumbnailFilename string `json:"thumbnail_filename"`
	VideoFilename     string `json:"video_filename"`
}

func GenerateVideoS3Key(info VideoInfo) string {
	extension := filepath.Ext(info.VideoFilename)
	return fmt.Sprintf("service/%s/%s/%s/video_segment/%s",
		info.UploadedBy,
		info.ServiceId,
		info.VideoId,
		fmt.Sprintf("%s%s", info.VideoId, extension),
	)
}

func GenerateVideoThumbnailS3Key(info VideoInfo) string {
	thumbnailFilename := generateThumbnailFilename(info.ThumbnailFilename)
	return fmt.Sprintf("service/%s/%s/%s/thumbnail/%s",
		info.UploadedBy,
		info.ServiceId,
		info.VideoId,
		thumbnailFilename,
	)
}

func generateThumbnailFilename(filename string) string {
	extension := filepath.Ext(filename)
	return "thumbnail" + extension
}

type ServiceInfo struct {
	UploadedBy string `json:"uploaded_by"`
	ServiceId  string `json:"service_id"`
	Filename   string `json:"filename"`
}

func GenerateServiceThumbnaiS3Key(info ServiceInfo) string {
	return fmt.Sprintf("service/%s/%s/%s",
		info.UploadedBy,
		info.ServiceId,
		utils.RenameFile(info.Filename, fmt.Sprintf("thumbnail_%s", info.ServiceId)),
	)
}

func GenerateUserProfilePictureS3Key(userId string, filename string) string {
	return fmt.Sprintf(
		"user/profile_picture/%s",
		utils.RenameFile(filename, fmt.Sprintf("%s_%s", uuid.NewString(), userId)),
	)
}

func GenerateCategoryImageS3Key(cateId string, filename string) string {
	return fmt.Sprintf(
		"category/image/%s",
		utils.RenameFile(filename, fmt.Sprintf("%s_%s", cateId, filename)),
	)
}

func GenerateServiceImageS3Key(serviceId string, filename string) string {
	return fmt.Sprintf(
		"service/%s/image/%s",
		serviceId, utils.RenameFile(filename, fmt.Sprintf("%s_%s", serviceId, filename)),
	)
}

func GenerateGroupProviderImageS3Key(groupProviderId string, filename string) string {
	return fmt.Sprintf(
		"group-provider/%s/image/%s",
		groupProviderId, utils.RenameFile(filename, fmt.Sprintf("%s_%s", groupProviderId, filename)),
	)
}

func GenerateKYCImageS3Key(kycImageId string, filename string) string {
	return fmt.Sprintf(
		"kyc/%s/image/%s",
		kycImageId, utils.RenameFile(filename, fmt.Sprintf("%s_%s", kycImageId, filename)),
	)
}

func GenerateCouponImageS3Key(couponId string, filename string) string {
	return fmt.Sprintf(
		"coupon/image/%s",
		utils.RenameFile(filename, fmt.Sprintf("%s_%s", couponId, filename)),
	)
}

func AddPublicCloudFrontDomain(s3Key string) string {
	return fmt.Sprintf("%s/%s", appconst.AWSCloudFrontPublicFile, s3Key)
}
