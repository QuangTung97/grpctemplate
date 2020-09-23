package backend

import (
	"context"
	domain "grpctemplate/domain/backend"
	rpc "grpctemplate/rpc/backend/v1"
	"time"

	"grpctemplate/domain/errors"

	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

// Service for gRPC
type Service struct {
	port domain.IPort
}

var _ rpc.HelloServer = &Service{}

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
	ctxzap.Extract(ctx).Info("context.Done", zap.Error(ctx.Err()))

	return &rpc.HelloResponse{}, nil
}
