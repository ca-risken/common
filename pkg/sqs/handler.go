package sqs

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/go-sqs-poller/worker/v5"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Handler is common interface for handling sqs message in the risken.
type Handler interface {
	HandleMessage(ctx context.Context, msg *types.Message) error
}

type HandlerFunc func(ctx context.Context, msg *types.Message) error

func (f HandlerFunc) HandleMessage(ctx context.Context, msg *types.Message) error {
	return f(ctx, msg)
}

// NonRetryableError indicates an error that the message cause the error cannot be retried.
type NonRetryableError struct {
	error
}

func WrapNonRetryable(err error) NonRetryableError {
	return NonRetryableError{err}
}

func (e NonRetryableError) Error() string {
	return fmt.Sprintf("NonRetryableError caused: %s", e.error.Error())
}

// InitializeHandler returns go-sqs-poller worker Handler from common Handler interface.
func InitializeHandler(h Handler) worker.Handler {
	return worker.HandlerFunc(func(msg *types.Message) error {
		ctx := context.Background()
		return h.HandleMessage(ctx, msg)
	})
}

func StatusLoggingHandler(logger logging.Logger, h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		now := time.Now()
		err := h.HandleMessage(ctx, msg)
		elapsed := time.Since(now).Seconds()
		items := map[string]interface{}{
			"elapsed": elapsed,
		}
		if err != nil {
			logger.WithItemsf(ctx, logging.WarnLevel, items, "handling message failed. err: %+v", err)
		} else {
			logger.WithItems(ctx, logging.InfoLevel, items, "handling message succeeded.")
		}
		return err
	})
}

// RetryableErrorHandler returns the Handler that returns nil when NonRetryableError occurred.
func RetryableErrorHandler(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		err := h.HandleMessage(ctx, msg)
		var target NonRetryableError
		if err != nil && !errors.As(err, &target) {
			return err
		}
		return nil
	})
}

func TracingHandler(serviceName string, h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *types.Message) error {
		span, tctx := ddtracer.StartSpanFromContext(ctx, serviceName)
		// TODO inherit trace from message
		err := h.HandleMessage(tctx, msg)
		span.Finish(ddtracer.WithError(err))
		return err
	})
}
