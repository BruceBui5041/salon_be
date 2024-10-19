package messagemodel

import (
	pb "salon_be/proto/salon_be/salon_be"
)

type VideoProcessStateInfo struct {
	VideoID           string               `json:"video_id"`
	ServiceID         string               `json:"service_id"`
	Timestamp         int64                `json:"timestamp"`
	Progress          int32                `json:"progress"`
	State             pb.ProcessState      `json:"state"`
	ErrorMsg          string               `json:"error_msg"`
	RequestResolution pb.ProcessResolution `json:"request_resolution"`
}
