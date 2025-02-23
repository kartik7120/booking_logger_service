package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

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

	pb.RegisterLoggerServiceServer(grpcServer, &grpcserver.LoggerServiceServer{
		LogRequest: &pb.LogRequest{},
	})

	// Starting the grpc server

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}

	// Gracefully shutting down the grpc server

	<-signalChan
	grpcServer.GracefulStop()
}
