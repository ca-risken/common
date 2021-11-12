package rpc

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func LoggingUnaryServerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		res, err := handler(ctx, req)
		method := info.FullMethod
		code := status.Code(err)
		elapsed := time.Since(start).Seconds()
		entry := logger.WithField("method", method).
			WithField("code", code.String()).
			WithField("elapsed", elapsed)
		if err != nil {
			entry = entry.WithField("error", shorten(err.Error()))
		}
		entry.Infof("receive from client")
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
