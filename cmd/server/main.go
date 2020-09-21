package main

import (
	"fmt"
	"grpctemplate/domain/backend"
	backend_rpc "grpctemplate/rpc/backend/v1"
	backend_service "grpctemplate/service/backend"
	"grpctemplate/service/interceptors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.DomainErrorUnaryInterceptor(),
		),
	)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	port := backend.NewPort()

	s := backend_service.NewService(port)

	backend_rpc.RegisterHelloService(server, &backend_rpc.HelloService{
		Hello: s.Hello,
	})

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		err = server.Serve(lis)
		if err != nil {
			log.Println("[ERROR]", err)
		}
	}()

	signal := <-exit
	fmt.Println("SIGNAL", signal)
	server.GracefulStop()
	wg.Wait()

	fmt.Println("Stop successfully")
}
