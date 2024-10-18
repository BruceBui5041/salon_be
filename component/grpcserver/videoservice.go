package grpcserver

import (
	pb "salon_be/proto/video_service/video_service" // import the generated protobuf package
)

type VideoServiceServer struct {
	pb.UnimplementedVideoProcessingServiceServer
}
