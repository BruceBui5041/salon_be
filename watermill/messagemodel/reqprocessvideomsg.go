package messagemodel

import (
	pb "salon_be/proto/video_service/video_service"
)

// RequestProcessVideoInfo represents the information about a newly uploaded video
type RequestProcessVideoInfo struct {
	Timestamp         string                `json:"timestamp"`
	RawVidS3Key       string                `json:"s3key"`
	UploadedBy        string                `json:"uploaded_by"`
	CourseId          string                `json:"course_id"`
	VideoId           string                `json:"video_id"`
	Retry             uint                  `json:"retry"`
	RequestResolution *pb.ProcessResolution `json:"request_resolution"`
}
