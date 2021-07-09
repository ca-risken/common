package xray

import (
	"context"
	"fmt"

	"github.com/aws/aws-xray-sdk-go/xray"
	"google.golang.org/grpc"
)

func AnnotateEnvTracingUnaryServerInterceptor(env string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := xray.AddAnnotation(ctx, "env", env); err != nil {
			// TODO logger
			//appLogger.Warnf("failed to annotate environment to x-ray: %+v", err)
			fmt.Printf("failed to annotate environment to x-ray: %+v", err)
		}
		return handler(ctx, req)
	}
}
