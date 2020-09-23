package backend

import (
	"context"
	"fmt"
	domain "grpctemplate/domain/backend"
	rpc "grpctemplate/rpc/backend/v1"
	"time"

	"grpctemplate/domain/errors"

	"github.com/golang/protobuf/ptypes"
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
	createdAt, err := ptypes.Timestamp(req.CreatedAt)
	if err != nil {
		return nil, errors.ErrInvalidTime
	}

	if req.Type < 1 || req.Type > 3 {
		return nil, errors.ErrInvalidCampaignType
	}

	input := domain.HelloInput{
		Type:      domain.CampaignType(req.Type),
		CreatedAt: createdAt,
	}

	err = s.port.Hello(input)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	<-ctx.Done()
	fmt.Println(ctx.Err())

	return &rpc.HelloResponse{}, nil
}
