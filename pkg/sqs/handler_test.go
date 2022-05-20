package sqs

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestRetryableErrorHandler(t *testing.T) {
	normalCase := HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		return nil
	})
	actual := RetryableErrorHandler(normalCase).HandleMessage(context.Background(), nil)
	assert.Nil(t, actual)

	retryableErrorCase := HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		return fmt.Errorf("some retryable error")
	})
	actual = RetryableErrorHandler(retryableErrorCase).HandleMessage(context.Background(), nil)
	assert.Error(t, actual)

	nonRetryableErrorCase := HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		err := fmt.Errorf("some non retryable error")
		return WrapNonRetryable(err)
	})
	actual = RetryableErrorHandler(nonRetryableErrorCase).HandleMessage(context.Background(), nil)
	assert.Nil(t, actual)
}

func TestTracingHandler(t *testing.T) {
	var tctx context.Context
	startTraceCase := HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		tctx = ctx
		return nil
	})
	actual := TracingHandler("test", startTraceCase).HandleMessage(context.Background(), nil)
	assert.NoError(t, actual)
	_, ok := ddtracer.SpanFromContext(tctx)
	assert.True(t, ok)

	errorCase := HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		tctx = ctx
		return fmt.Errorf("some error")
	})
	actual = TracingHandler("test", errorCase).HandleMessage(context.Background(), nil)
	assert.Error(t, actual)
	_, ok = ddtracer.SpanFromContext(tctx)
	assert.True(t, ok)
}
