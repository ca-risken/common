package xray

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-xray-sdk-go/header"
	"github.com/aws/aws-xray-sdk-go/xray"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
)

type QueueMessage struct {
	ScanOnly bool `json:"scan_only,string"`
}

func MessageTracingHandler(env, segmentName string, next mimosasqs.Handler) mimosasqs.Handler {
	return mimosasqs.HandlerFunc(func(ctx context.Context, msg *sqs.Message) error {
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
