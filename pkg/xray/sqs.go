package xray

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gassara-kys/go-sqs-poller/worker/v4"
)

type QueueMessage struct {
	ScanOnly bool `json:"scan_only,string"`
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, message *sqs.Message) error
}

func MessageTracingHandler(env, segmentName string, next MessageHandler) worker.Handler {
	return worker.HandlerFunc(func(msg *sqs.Message) error {
		ctx := context.Background()
		ctx, segment := xray.BeginSegment(ctx, segmentName)

		inheritTrace := true
		mb := &QueueMessage{}
		err := json.Unmarshal([]byte(aws.StringValue(msg.Body)), mb)
		if err != nil || mb.ScanOnly {
			inheritTrace = false
		}

		var th *header.Header
		thString, ok := msg.Attributes["AWSTraceHeader"]
		if inheritTrace && ok {
			th = header.FromString(aws.StringValue(thString))
			segment.TraceID = th.TraceID
			segment.ParentID = th.ParentID
			segment.Sampled = th.SamplingDecision == header.Sampled
		}
		if err := xray.AddAnnotation(ctx, "env", env); err != nil {
			// TODO logger
			fmt.Printf("failed to annotate environment to x-ray: %+v", err)
		}

		err = next.HandleMessage(ctx, msg)
		segment.Close(err)
		return err
	})
}
