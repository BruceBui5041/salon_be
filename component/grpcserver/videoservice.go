package grpcserver

import (
	pb "salon_be/proto/salon_be/salon_be" // import the generated protobuf package
)

type VideoServiceServer struct {
	pb.UnimplementedVideoProcessingServiceServer
}
