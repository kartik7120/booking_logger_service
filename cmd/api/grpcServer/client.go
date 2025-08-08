package grpcserver

import (
	"context"

	pb "github.com/booking_logger_service/cmd/api"
	"google.golang.org/grpc"
)

func NewLoggerServiceClient() pb.LoggerServiceClient {

	conn, err := grpc.NewClient("localhost:1100")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewLoggerServiceClient(conn)

	return client
}

func PrintLog(client pb.LoggerServiceClient, log *pb.LogRequest) {

	_, err := client.Log(context.Background(), log)

	if err != nil {
		panic(err)
	}
}
