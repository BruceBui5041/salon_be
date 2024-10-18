package grpcserver

import (
	pb "video_server/proto/video_service/video_service" // import the generated protobuf package
)

type VideoServiceServer struct {
	pb.UnimplementedVideoProcessingServiceServer
}
