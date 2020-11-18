package main

import (
	"context"
	"fmt"
	backend_rpc "grpctemplate/rpc/backend/v1"
	lib_rpc "grpctemplate/rpc/lib/v1"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
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
		st, ok := status.FromError(err)
		if ok {
			fmt.Println("Code:", st.Code())
			fmt.Println("Message:", st.Message())
		}
		panic(err)
	}
}
