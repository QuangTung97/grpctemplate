package tracing

import (
	"context"

	"github.com/google/uuid"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func getTraceIdFromContext(ctx context.Context) (uuid.UUID, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("trace.request.id")
		if len(values) > 0 {
			header := values[0]
			id, err := uuid.Parse(header)
			if err == nil {
				return id, true
			}
		}
	}
	return uuid.UUID{}, false
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		id, ok := getTraceIdFromContext(ctx)
		if !ok {
			id = uuid.New()
		}

		tags := grpc_ctxtags.Extract(ctx)
		tags = tags.Set("trace.request.id", id)
		ctx = grpc_ctxtags.SetInContext(ctx, tags)

		return handler(ctx, req)
	}
}
