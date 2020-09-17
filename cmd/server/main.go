package main

import (
	"grpctemplate/rpc/backend/v1"
	"grpctemplate/service/backend"
	"net"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()

	s := backend_service.NewService()

	service := &backend.HelloService{
		Hello: s.Hello,
	}

	backend.RegisterHelloService(server, service)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)
	panic(err)
}
