package sqs

import (
	"context"
	"fmt"
	"testing"

	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func TestRetryableErrorHandler(t *testing.T) {
	normalCase := HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		return nil
	})
	actual := RetryableErrorHandler(normalCase).HandleMessage(context.Background(), nil)
	assert.Nil(t, actual)

	retryableErrorCase := HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		return fmt.Errorf("some retryable error")
	})
	actual = RetryableErrorHandler(retryableErrorCase).HandleMessage(context.Background(), nil)
	assert.Error(t, actual)

	nonRetryableErrorCase := HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		err := fmt.Errorf("some non retryable error")
		return WrapNonRetryable(err)
	})
	actual = RetryableErrorHandler(nonRetryableErrorCase).HandleMessage(context.Background(), nil)
	assert.Nil(t, actual)
}

func TestTracingHandler(t *testing.T) {
	startTraceCase := HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		_, ok := ddtracer.SpanFromContext(ctx)
		assert.True(t, ok)
		return nil
	})
	actual := TracingHandler("test", startTraceCase).HandleMessage(context.Background(), nil)
	assert.NoError(t, actual)
}
