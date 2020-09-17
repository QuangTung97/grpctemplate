package main

import (
	"context"
	backend_rpc "grpctemplate/rpc/backend/v1"
	lib_rpc "grpctemplate/rpc/lib/v1"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := backend_rpc.NewHelloClient(conn)

	ctx := context.Background()

	createdAt, err := ptypes.TimestampProto(time.Now())
	if err != nil {
		panic(err)
	}

	req := &backend_rpc.HelloRequest{
		Type:      lib_rpc.CampaignType_CAMPAIGN_TYPE_MERCHANT,
		CreatedAt: createdAt,
	}

	_, err = client.Hello(ctx, req)
	if err != nil {
		panic(err)
	}
}
