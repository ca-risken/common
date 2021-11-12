package rpc

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const longErrorMessage = "あいうえおかきくけこさしすせそたちつてとなにぬねのあいうえおかきくけこさしすせそたちつてとなにぬねのあいうえおかきくけこさしすせそたちつてとなにぬねのあいうえおかきくけこさしすせそたちつてとなにぬねのあ"

func TestLoggingUnaryServerInterceptor(t *testing.T) {
	logger := logrus.New()
	interceptor := LoggingUnaryServerInterceptor(logger)
	info := &grpc.UnaryServerInfo{
		Server:     nil,
		FullMethod: "/some/method",
	}
	// error case
	res, err := interceptor(context.Background(), "request", info,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return "response", errors.New(longErrorMessage)
		})
	assert.Error(t, err)
	assert.Equal(t, "response", res)
	// normal case
	res, err = interceptor(context.Background(), "request", info,
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return "response", nil
		})
	assert.NoError(t, err)
	assert.Equal(t, "response", res)
}
