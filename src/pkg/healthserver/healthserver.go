package healthserver

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type HealthServer struct{}

func NewHealthServer() grpc_health_v1.HealthServer {
	return &HealthServer{}
}

func (h *HealthServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

func (h *HealthServer) Check(context.Context, *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *HealthServer) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}
