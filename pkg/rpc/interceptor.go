package rpc

import (
	"context"
	"time"

	"github.com/ca-risken/common/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryServerInterceptor(logger logging.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		res, err := handler(ctx, req)
		method := info.FullMethod
		code := status.Code(err)
		elapsed := time.Since(start).Seconds()
		items := map[string]interface{}{
			"method":  method,
			"code":    code.String(),
			"elapsed": elapsed,
		}
		if err != nil {
			items["error"] = shorten(err.Error())
		}
		logger.WithItems(ctx, logging.InfoLevel, items, "receive from client")
		return res, err
	}
}

func shorten(msg string) string {
	maxLength := 100
	if len(msg) <= maxLength {
		return msg
	}
	return string([]rune(msg)[:maxLength]) + "..."
}
