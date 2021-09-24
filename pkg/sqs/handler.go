package sqs

import (
	"context"

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

// InitializeHandler returns go-sqs-poller worker Handler from common Handler interface.
func InitializeHandler(h Handler) worker.Handler {
	return worker.HandlerFunc(func(msg *awssqs.Message) error {
		ctx := context.Background()
		return h.HandleMessage(ctx, msg)
	})
}

func StatusLoggingHandler(logger *logrus.Logger, h Handler) Handler {
	return HandlerFunc(func(ctx context.Context, msg *awssqs.Message) error {
		err := h.HandleMessage(ctx, msg)
		if err != nil {
			logger.Warnf("handling message failed. err: %+v", err)
		} else {
			logger.Infof("handling message succeeded.")
		}
		return err
	})
}
