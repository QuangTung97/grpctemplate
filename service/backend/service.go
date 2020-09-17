package backend

import (
	"context"
	domain "grpctemplate/domain/backend"
	rpc "grpctemplate/rpc/backend/v1"
)

// Service for gRPC
type Service struct {
	port domain.IPort
}

// NewService create a new Service
func NewService(port domain.IPort) *Service {
	return &Service{
		port: port,
	}
}

// Hello do hello
func (s *Service) Hello(ctx context.Context, req *rpc.HelloRequest,
) (*rpc.HelloResponse, error) {
	return &rpc.HelloResponse{}, nil
}
