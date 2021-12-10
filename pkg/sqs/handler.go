package sqs

import (
	"context"
	"errors"
	"fmt"
	"time"

	awssqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gassara-kys/go-sqs-poller/worker/v4"
	"github.com/sirupsen/logrus"
)

// Handler is common interface for handling sqs message in the risken.
type Handler interface {
	HandleMessage(ctx context.Context, msg *awssqs.Message) error
}

type HandlerFunc func(ctx context.Context, msg *awssqs.Message) error

func (f HandlerFunc) HandleMessage(ctx context.Context, msg *awssqs.Message) error {
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
	return worker.HandlerFunc(func(msg *awssqs.Message) error {
		ctx := context.Background()
		return h.HandleMessage(ctx, msg)
	})
}

func StatusLoggingHandler(logger *logrus.Logger, h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		now := time.Now()
		err := h.HandleMessage(ctx, msg)
		elapsed := time.Now().Sub(now).Seconds()
		slogger := logger.WithField("elapsed", elapsed)
		if err != nil {
			slogger.Warnf("handling message failed. err: %+v", err)
		} else {
			slogger.Infof("handling message succeeded.")
		}
		return err
	})
}

// RetryableErrorHandler returns the Handler that returns nil when NonRetryableError occurred.
func RetryableErrorHandler(h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		err := h.HandleMessage(ctx, msg)
		var target NonRetryableError
		if err != nil && !errors.As(err, &target) {
			return err
		}
		return nil
	})
}
