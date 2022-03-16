package sqs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/ca-risken/common/pkg/logging"
	"github.com/ca-risken/core/proto/finding"
)

type Finalizer struct {
	findingClient finding.FindingServiceClient
}

func NewFinalizer(findingSvcAddr string) *Finalizer {
	return &Finalizer{
		findingClient: newFindingClient(findingSvcAddr),
	}
}

// DataSourceRecommnend needs to implement scan failure recommendations contents per datasource
type DataSourceRecommnend interface {
	// DataSource returns the DataSource identifier string.(e.g. 'aws:portscan')
	DataSource() string
	// ScanFailureRisk returns risk information in case of scan failure.
	ScanFailureRisk() string
	// ScanFailureRecommend returns information on what action is required in case of scan failure.
	ScanFailureRecommend() string
}

// FinalizeHandler returns a Handler that wraps the termination process
func (f *Finalizer) FinalizeHandler(datasource DataSourceRecommnend, next Handler) Handler {
	return HandlerFunc(func(ctx context.Context, sqsMsg *sqs.Message) error {
		err := next.HandleMessage(ctx, sqsMsg)
		projectID, parseErr := parseProjectFromMessage(aws.StringValue(sqsMsg.Body))
		if err != nil {
			appLogger.Errorf("Invalid message(failed to get projetc_id): sqsMsg=%+v, err=%+v", sqsMsg, parseErr)
			return f.Final(ctx, nil, datasource, err)
		}
		return f.Final(ctx, &projectID, datasource, err)
	})
}

type Message struct {
	ProjectID uint32 `json:"project_id,omitempty"`
}

func parseProjectFromMessage(msg string) (uint32, error) {
	message := &Message{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return 0, err
	}
	return message.ProjectID, nil
}

type recommend struct {
	Risk           string `json:"risk,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

// Final summarizes the termination scan process
func (f *Finalizer) Final(ctx context.Context, projectID *uint32, datasource DataSourceRecommnend, err error) error {
	if projectID == nil {
		// Unknown project
		appLogger.Notifyf(logging.ErrorLevel, "Unknown project, err: %+v", err)
		return err
	}
	if err != nil {
		// Scan failed
		if putErr := f.putScanFinding(ctx, projectID, &ScanFinding{
			ProjectID:    *projectID,
			DataSource:   datasource.DataSource(),
			Status:       "Error",
			ErrorMessage: err.Error(),
			Recommend: recommend{
				Risk:           datasource.ScanFailureRisk(),
				Recommendation: datasource.ScanFailureRecommend(),
			},
		}); putErr != nil {
			appLogger.Notifyf(logging.ErrorLevel, "Failed to putScanFinding (scan failed), project_id: %d, err: %+v", *projectID, putErr)
			return err
		}
		return err
	}

	// Scan succeeded
	if putErr := f.putScanFinding(ctx, projectID, &ScanFinding{
		ProjectID:  *projectID,
		DataSource: datasource.DataSource(),
		Status:     "OK",
		Recommend: recommend{
			Risk:           datasource.ScanFailureRisk(),
			Recommendation: datasource.ScanFailureRecommend(),
		},
	}); putErr != nil {
		appLogger.Notifyf(logging.ErrorLevel, "Failed to putScanFinding (scan succeeded), project_id: %d, err: %+v", *projectID, putErr)
		return nil
	}
	return nil
}

type ScanFinding struct {
	ProjectID    uint32    `json:"project_id,omitempty"`
	DataSource   string    `json:"data_source,omitempty"`
	Status       string    `json:"status,omitempty"`
	ErrorMessage string    `json:"error_message,omitempty"`
	Recommend    recommend `json:"recommend,omitempty"`
}

func (f *Finalizer) putScanFinding(ctx context.Context, projectID *uint32, sf *ScanFinding) error {
	if projectID == nil || sf == nil {
		return nil // nop
	}
	score := float32(0.0)
	desc := fmt.Sprintf("Successfully scanned %s", sf.DataSource)
	if sf.ErrorMessage != "" {
		desc = fmt.Sprintf("Failed to scan %s", sf.DataSource)
		score = 0.8
	}

	buf, err := json.Marshal(sf)
	if err != nil {
		return err
	}
	// PutFinding
	resp, err := f.findingClient.PutFinding(ctx, &finding.PutFindingRequest{
		Finding: &finding.FindingForUpsert{
			Description:      desc,
			DataSource:       "RISKEN",
			DataSourceId:     fmt.Sprintf("%s-scan-status", sf.DataSource),
			ResourceName:     sf.DataSource,
			ProjectId:        sf.ProjectID,
			OriginalScore:    score,
			OriginalMaxScore: 1.0,
			Data:             string(buf),
		},
	})
	if err != nil {
		return err
	}
	if _, err := f.findingClient.PutRecommend(ctx, &finding.PutRecommendRequest{
		ProjectId:      sf.ProjectID,
		FindingId:      resp.Finding.FindingId,
		DataSource:     "RISKEN",
		Type:           fmt.Sprintf("ScanError/%s", sf.DataSource),
		Risk:           sf.Recommend.Risk,
		Recommendation: sf.Recommend.Recommendation,
	}); err != nil {
		return fmt.Errorf("Failed to put scan finding recommned, finding_id=%d, error=%+v", resp.Finding.FindingId, err)
	}
	return nil
}
