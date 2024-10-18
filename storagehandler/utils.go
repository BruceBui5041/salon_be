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
	CourseId          string `json:"course_id"`
	VideoId           string `json:"video_id"`
	ThumbnailFilename string `json:"thumbnail_filename"`
	VideoFilename     string `json:"video_filename"`
}

func GenerateVideoS3Key(info VideoInfo) string {
	extension := filepath.Ext(info.VideoFilename)
	return fmt.Sprintf("course/%s/%s/%s/video_segment/%s",
		info.UploadedBy,
		info.CourseId,
		info.VideoId,
		fmt.Sprintf("%s%s", info.VideoId, extension),
	)
}

func GenerateVideoThumbnailS3Key(info VideoInfo) string {
	thumbnailFilename := generateThumbnailFilename(info.ThumbnailFilename)
	return fmt.Sprintf("course/%s/%s/%s/thumbnail/%s",
		info.UploadedBy,
		info.CourseId,
		info.VideoId,
		thumbnailFilename,
	)
}

func generateThumbnailFilename(filename string) string {
	extension := filepath.Ext(filename)
	return "thumbnail" + extension
}

type CourseInfo struct {
	UploadedBy string `json:"uploaded_by"`
	CourseId   string `json:"course_id"`
	Filename   string `json:"filename"`
}

func GenerateCourseThumbnaiS3Key(info CourseInfo) string {
	return fmt.Sprintf("course/%s/%s/%s",
		info.UploadedBy,
		info.CourseId,
		utils.RenameFile(info.Filename, fmt.Sprintf("thumbnail_%s", info.CourseId)),
	)
}

func GenerateUserProfilePictureS3Key(userId string, filename string) string {
	return fmt.Sprintf(
		"user/profile_picture/%s",
		utils.RenameFile(filename, fmt.Sprintf("%s_%s", uuid.NewString(), userId)),
	)
}

func AddPublicCloudFrontDomain(s3Key string) string {
	return fmt.Sprintf("%s/%s", appconst.AWSCloudFrontPublicFile, s3Key)
}
