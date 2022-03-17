package service

import (
	context "context"

	"google.golang.org/grpc"

	pb "google.golang.org/grpc/health/grpc_health_v1"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	resp := &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}

	return resp, nil
}
func (h *HealthHandler) Watch(req *pb.HealthCheckRequest, srv pb.Health_WatchServer) error {
	resp := &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}

	for {
		err := srv.Send(resp)
		if err != nil {
			return err
		}
	}
}

func (h *HealthHandler) registerHandler(s *grpc.Server) {
	pb.RegisterHealthServer(s, h)
}
