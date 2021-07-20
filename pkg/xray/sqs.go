package xray

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gassara-kys/go-sqs-poller/worker/v4"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, message *sqs.Message) error
}

func MessageTracingHandler(env, segmentName string, next MessageHandler) worker.Handler {
	return worker.HandlerFunc(func(msg *sqs.Message) error {
		ctx := context.Background()
		ctx, segment := xray.BeginSegment(ctx, segmentName)

		var th *header.Header
		thString, ok := msg.Attributes["AWSTraceHeader"]
		if ok {
			th = header.FromString(aws.StringValue(thString))
			segment.TraceID = th.TraceID
			segment.ParentID = th.ParentID
			segment.Sampled = th.SamplingDecision == header.Sampled
		}
		if err := xray.AddAnnotation(ctx, "env", env); err != nil {
			// TODO logger
			fmt.Printf("failed to annotate environment to x-ray: %+v", err)
		}

		err := next.HandleMessage(ctx, msg)
		segment.Close(err)
		return err
	})
}
