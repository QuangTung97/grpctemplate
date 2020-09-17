package backend_service

import (
	"context"
	rpc "grpctemplate/rpc/backend/v1"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Hello(ctx context.Context, req *rpc.HelloRequest,
) (*rpc.HelloResponse, error) {
	return &rpc.HelloResponse{}, nil
}
