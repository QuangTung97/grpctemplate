package main

import (
	"grpctemplate/domain/backend"
	backend_rpc "grpctemplate/rpc/backend/v1"
	backend_service "grpctemplate/service/backend"
	"grpctemplate/service/interceptors"
	"net"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.DomainErrorUnaryInterceptor(),
		),
	)

	port := backend.NewPort()

	s := backend_service.NewService(port)

	backend_rpc.RegisterHelloService(server, &backend_rpc.HelloService{
		Hello: s.Hello,
	})

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	err = server.Serve(lis)
	panic(err)
}
