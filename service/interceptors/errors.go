package interceptors

import (
	"context"
	"grpctemplate/domain/errors"
	"strconv"

	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func domainErrorCodeToGRPC(code string) codes.Code {
	num, err := strconv.ParseInt(code[:2], 10, 32)
	if err != nil {
		panic(err)
	}
	return codes.Code(num)
}

// DomainErrorUnaryInterceptor intercepts domain's errors
func DomainErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		tags := grpc_ctxtags.Extract(ctx)
		tags = tags.Set("user.id", 111222)
		ctx = grpc_ctxtags.SetInContext(ctx, tags)

		res, err := handler(ctx, req)

		if err != nil {
			domainErr, ok := err.(errors.Error)
			if !ok {
				st := status.New(codes.Internal, "Internal server error")
				return nil, st.Err()
			}

			code := domainErrorCodeToGRPC(domainErr.Code)
			st := status.New(code, domainErr.Message)

			return nil, st.Err()
		}

		return res, nil
	}
}
