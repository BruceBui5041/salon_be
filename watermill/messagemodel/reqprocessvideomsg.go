package messagemodel

import (
	pb "salon_be/proto/salon_be/salon_be"
)

// RequestProcessVideoInfo represents the information about a newly uploaded video
type RequestProcessVideoInfo struct {
	Timestamp         string                `json:"timestamp"`
	RawVidS3Key       string                `json:"s3key"`
	UploadedBy        string                `json:"uploaded_by"`
	ServiceId         string                `json:"service_id"`
	VideoId           string                `json:"video_id"`
	Retry             uint                  `json:"retry"`
	RequestResolution *pb.ProcessResolution `json:"request_resolution"`
}
