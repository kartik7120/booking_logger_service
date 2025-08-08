package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/booking_logger_service/cmd/api"
	grpcserver "github.com/booking_logger_service/cmd/api/grpcServer"
)

func main() {

	lis, err := net.Listen("tcp", ":1100")

	if err != nil {
		panic(err)
	}

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Creating a new grpc server

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// Registering the grpc server

	pb.RegisterLoggerServiceServer(grpcServer, &grpcserver.LoggerServiceServer{})

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Starting the grpc server

	go func() {
		fmt.Println("Logging service started")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// Gracefully shutting down the grpc server

	<-signalChan
	fmt.Println("Shutting down the logging service")
	grpcServer.GracefulStop()
}
