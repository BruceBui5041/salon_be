package messagemodel

import (
	pb "video_server/proto/video_service/video_service"
)

type VideoProcessStateInfo struct {
	VideoID           string               `json:"video_id"`
	CourseID          string               `json:"course_id"`
	Timestamp         int64                `json:"timestamp"`
	Progress          int32                `json:"progress"`
	State             pb.ProcessState      `json:"state"`
	ErrorMsg          string               `json:"error_msg"`
	RequestResolution pb.ProcessResolution `json:"request_resolution"`
}
