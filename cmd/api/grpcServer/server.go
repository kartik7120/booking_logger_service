package grpcserver

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/booking_logger_service/cmd/api"
	"github.com/booking_logger_service/cmd/logger"
)

type LoggerServiceServer struct {
	pb.UnimplementedLoggerServiceServer
	LogRequest *pb.LogRequest
	mu         sync.Mutex
}

func (s *LoggerServiceServer) Log(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	logger, err, provider := logger.GetLoggerWithContext(ctx)

	if err != nil {
		return &pb.LogResponse{
			Status:  "ERROR",
			Message: "Failed to create logger",
			Error:   err.Error(),
		}, err
	}

	if in.GetLevel() == pb.LogLevel_DEBUG {
		logger.Debug(in.GetMessage())
	} else if in.GetLevel() == pb.LogLevel_INFO {
		logger.Info(in.GetMessage())
	} else if in.GetLevel() == pb.LogLevel_WARNING {
		logger.Warn(in.GetMessage())
	} else if in.GetLevel() == pb.LogLevel_ERROR {
		logger.Error(in.GetMessage())
	} else if in.GetLevel() == pb.LogLevel_FATAL {
		logger.Fatal(in.GetMessage())
	}

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.Canceled {
				cancel()
				return &pb.LogResponse{
					Status:  "ERROR",
					Message: "Context was cancelled",
					Error:   "",
				}, nil
			}
			cancel()
			return &pb.LogResponse{
				Status:  "OK",
				Message: "Log message received",
				Error:   "",
			}, nil
		case <-signalChan:
			cancel()
			provider.Shutdown(ctx)
			os.Exit(0)
			return &pb.LogResponse{
				Status: "Interrupted",
				Error:  "Received an interrupt signal",
			}, nil
		}
	}

}
