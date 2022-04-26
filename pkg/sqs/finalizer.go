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
	datasource     string
	recommendation *DataSourceRecommnend
	findingClient  finding.FindingServiceClient
}

// NewFinalizer returns a Finalizer instance and generates recommendation contents for scan failure from the datasource and settingURL parameter.
// If set overrideRecommendation, the above contents will be overridden.
func NewFinalizer(datasource, settingURL, findingSvcAddr string, overrideRecommendation *DataSourceRecommnend) (*Finalizer, error) {
	r, err := generateRecommendation(datasource, settingURL, overrideRecommendation)
	if err != nil {
		return nil, err
	}
	fc, err := newFindingClient(findingSvcAddr)
	if err != nil {
		return nil, err
	}
	return &Finalizer{
		datasource:     datasource,
		recommendation: r,
		findingClient:  fc,
	}, nil
}

const (
	scanFailureRiskTemplate      = "Failed to scan %s, So you are not gathering the latest security threat information."
	scanFailureRecommendTemplate = `Please review the following items and rescan,
	- Ensure the error message of the DataSource.
	- Ensure the access rights you set for the DataSource and the reachability of the network.
	- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
	- %s
	- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`
)

type DataSourceRecommnend struct {
	// ScanFailureRisk is risk information in case of scan failure.
	ScanFailureRisk string `json:"scan_failure_risk,omitempty"`
	// ScanFailureRecommendation is information on what action is required in case of scan failure.
	ScanFailureRecommendation string `json:"scan_failure_recommendation,omitempty"`
}

func generateRecommendation(datasource, settingURL string, override *DataSourceRecommnend) (*DataSourceRecommnend, error) {
	if datasource == "" {
		return nil, fmt.Errorf("Required datasource")
	}
	if settingURL == "" && override == nil {
		return nil, fmt.Errorf("Required settingURL")
	}
	recommend := DataSourceRecommnend{
		ScanFailureRisk:           fmt.Sprintf(scanFailureRiskTemplate, datasource),
		ScanFailureRecommendation: fmt.Sprintf(scanFailureRecommendTemplate, settingURL),
	}
	if override != nil {
		recommend.ScanFailureRisk = override.ScanFailureRisk
		recommend.ScanFailureRecommendation = override.ScanFailureRecommendation
	}
	return &recommend, nil
}

// FinalizeHandler returns a Handler that wraps the termination process
func (f *Finalizer) FinalizeHandler(next Handler) Handler {
	return HandlerFunc(func(ctx context.Context, sqsMsg *sqs.Message) error {
		err := next.HandleMessage(ctx, sqsMsg)
		projectID, parseErr := parseProjectFromMessage(aws.StringValue(sqsMsg.Body))
		if parseErr != nil {
			appLogger.Errorf("Invalid message(failed to get project_id): sqsMsg=%+v, err=%+v", sqsMsg, parseErr)
			return f.Final(ctx, nil, err)
		}
		return f.Final(ctx, &projectID, err)
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

// Final summarizes the termination scan process
func (f *Finalizer) Final(ctx context.Context, projectID *uint32, err error) error {
	if projectID == nil {
		// Unknown project
		appLogger.Notifyf(logging.ErrorLevel, "Unknown project, err: %+v", err)
		return err
	}
	if err != nil {
		// Scan failed
		if putErr := f.putScanFinding(ctx, projectID, &ScanFinding{
			ProjectID:    *projectID,
			DataSource:   f.datasource,
			Status:       "Error",
			ErrorMessage: err.Error(),
			Recommendation: &DataSourceRecommnend{
				ScanFailureRisk:           f.recommendation.ScanFailureRisk,
				ScanFailureRecommendation: f.recommendation.ScanFailureRecommendation,
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
		DataSource: f.datasource,
		Status:     "OK",
		Recommendation: &DataSourceRecommnend{
			ScanFailureRisk:           f.recommendation.ScanFailureRisk,
			ScanFailureRecommendation: f.recommendation.ScanFailureRecommendation,
		},
	}); putErr != nil {
		appLogger.Notifyf(logging.ErrorLevel, "Failed to putScanFinding (scan succeeded), project_id: %d, err: %+v", *projectID, putErr)
		return nil
	}
	return nil
}

type ScanFinding struct {
	ProjectID      uint32                `json:"project_id,omitempty"`
	DataSource     string                `json:"data_source,omitempty"`
	Status         string                `json:"status,omitempty"`
	ErrorMessage   string                `json:"error_message,omitempty"`
	Recommendation *DataSourceRecommnend `json:"recommendation,omitempty"`
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
		Risk:           sf.Recommendation.ScanFailureRisk,
		Recommendation: sf.Recommendation.ScanFailureRecommendation,
	}); err != nil {
		return fmt.Errorf("Failed to put scan finding recommned, finding_id=%d, error=%+v", resp.Finding.FindingId, err)
	}
	return nil
}
