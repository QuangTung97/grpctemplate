package main

import (
	"context"
	"fmt"
	"grpctemplate/domain/backend"
	lib "grpctemplate/lib"
	backend_rpc "grpctemplate/rpc/backend/v1"
	backend_service "grpctemplate/service/backend"
	"grpctemplate/service/interceptors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"google.golang.org/grpc"
)

func initServer(logger *zap.Logger) *grpc.Server {
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	decider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return true
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			lib.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_zap.PayloadUnaryServerInterceptor(logger, decider),
			interceptors.DomainErrorUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			grpc_zap.PayloadStreamServerInterceptor(logger, decider),
		),
	)

	port := backend.NewPort()

	s := backend_service.NewService(port)

	backend_rpc.RegisterHelloServer(server, s)

	grpc_prometheus.Register(server)
	grpc_prometheus.EnableHandlingTimeHistogram()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true}),
	)
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	backend_rpc.RegisterHelloHandlerFromEndpoint(ctx, mux, "localhost:5000", opts)

	http.Handle("/api/", mux)
	http.Handle("/metrics", promhttp.Handler())

	return server
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	server := initServer(logger)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err = server.Serve(lis)
		if err != nil {
			logger.Error("serve", zap.Error(err))
		}
	}()

	httpServer := http.Server{
		Addr: ":5123",
	}

	go func() {
		defer wg.Done()

		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			return
		}
		if err != nil {
			logger.Error("httpServer Listen", zap.Error(err))
		}
	}()

	signal := <-exit
	fmt.Println("SIGNAL", signal)

	server.GracefulStop()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error("httpServer Shutdown", zap.Error(err))
	}

	wg.Wait()

	fmt.Println("Stop successfully")
}
