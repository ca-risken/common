package sqs

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/finding/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	sampleProjectID  uint32 = 1
	sampleDataSource string = "namespace:datasource"
	sampleSettingURL string = "https://docs.security-hub.jp/"
)

func TestGenerateRecommendation(t *testing.T) {
	type args struct {
		datasource string
		settingURL string
		override   *DataSourceRecommnend
	}
	cases := []struct {
		name    string
		input   args
		want    *DataSourceRecommnend
		wantErr bool
	}{
		{
			name: "OK",
			input: args{
				datasource: sampleDataSource,
				settingURL: sampleSettingURL,
				override:   nil,
			},
			want: &DataSourceRecommnend{
				ScanFailureRisk: "Failed to scan namespace:datasource, So you are not gathering the latest security threat information.",
				ScanFailureRecommendation: `Please review the following items and rescan,
	- Ensure the error message of the DataSource.
	- Ensure the access rights you set for the DataSource and the reachability of the network.
	- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
	- https://docs.security-hub.jp/
	- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
			},
		},
		{
			name: "OK Override",
			input: args{
				datasource: sampleDataSource,
				settingURL: "",
				override: &DataSourceRecommnend{
					ScanFailureRisk:           "overriden risk",
					ScanFailureRecommendation: "overriden recommendation",
				},
			},
			want: &DataSourceRecommnend{
				ScanFailureRisk:           "overriden risk",
				ScanFailureRecommendation: "overriden recommendation",
			},
		},
		{
			name: "NG required datasource",
			input: args{
				datasource: "",
				settingURL: sampleSettingURL,
				override:   nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "NG required settingURL",
			input: args{
				datasource: sampleDataSource,
				settingURL: "",
				override:   nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := generateRecommendation(c.input.datasource, c.input.settingURL, c.input.override)
			if (c.wantErr && err == nil) || (!c.wantErr && err != nil) {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("Unexpected data match: want=%+v, got=%+v", c.want, got)
			}
		})
	}

}

func TestFinal(t *testing.T) {
	type args struct {
		ProjectID *uint32
		Err       error
	}
	mockClient := mocks.FindingServiceClient{}
	type mockResponse struct {
		PutFindingResp   *finding.PutFindingResponse
		PutFindingErr    error
		PutRecommendResp *finding.PutRecommendResponse
		PutRecommendErr  error
	}

	cases := []struct {
		name     string
		input    args
		mockResp mockResponse
		wantErr  bool
	}{
		{
			name: "OK(no scan error)",
			input: args{
				ProjectID: &sampleProjectID,
				Err:       nil,
			},
			mockResp: mockResponse{
				PutFindingResp:   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1}},
				PutFindingErr:    nil,
				PutRecommendResp: &finding.PutRecommendResponse{Recommend: &finding.Recommend{}},
				PutRecommendErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "OK(exists scan error)",
			input: args{
				ProjectID: &sampleProjectID,
				Err:       errors.New("Failed to scan"),
			},
			mockResp: mockResponse{
				PutFindingResp:   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1}},
				PutFindingErr:    nil,
				PutRecommendResp: &finding.PutRecommendResponse{Recommend: &finding.Recommend{}},
				PutRecommendErr:  nil,
			},
			wantErr: true,
		},
		{
			name: "ProjectID is nil(error)",
			input: args{
				ProjectID: nil,
				Err:       errors.New("Failed to scan"),
			},
			mockResp: mockResponse{},
			wantErr:  true,
		},
		{
			name: "ProjectID is nil(no error)",
			input: args{
				ProjectID: nil,
				Err:       nil,
			},
			mockResp: mockResponse{},
			wantErr:  false,
		},
		{
			name: "Failed to PutFinding API",
			input: args{
				ProjectID: &sampleProjectID,
				Err:       nil,
			},
			mockResp: mockResponse{
				PutFindingResp:   nil,
				PutFindingErr:    errors.New("Failed to PutFinding"),
				PutRecommendResp: nil,
				PutRecommendErr:  nil,
			},
			wantErr: false,
		},
		{
			name: "Failed to PutRecommend API",
			input: args{
				ProjectID: &sampleProjectID,
				Err:       nil,
			},
			mockResp: mockResponse{
				PutFindingResp:   &finding.PutFindingResponse{Finding: &finding.Finding{FindingId: 1}},
				PutFindingErr:    nil,
				PutRecommendResp: nil,
				PutRecommendErr:  errors.New("Failed to PutRecommend"),
			},
			wantErr: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// setting mock client
			f := Finalizer{
				datasource: sampleDataSource,
				recommendation: &DataSourceRecommnend{
					ScanFailureRisk:           fmt.Sprintf(scanFailureRiskTemplate, sampleDataSource),
					ScanFailureRecommendation: fmt.Sprintf(scanFailureRecommendTemplate, sampleSettingURL),
				},
				findingClient: &mockClient,
			}
			if c.mockResp.PutFindingResp != nil || c.mockResp.PutFindingErr != nil {
				mockClient.On("PutFinding", mock.Anything, mock.Anything).Return(
					c.mockResp.PutFindingResp, c.mockResp.PutFindingErr).Once()
			}
			if c.mockResp.PutRecommendResp != nil || c.mockResp.PutRecommendErr != nil {
				mockClient.On("PutRecommend", mock.Anything, mock.Anything).Return(
					c.mockResp.PutRecommendResp, c.mockResp.PutRecommendErr).Once()
			}

			// exec function
			ctx := context.Background()
			err := f.Final(ctx, c.input.ProjectID, c.input.Err)
			if err == nil && c.wantErr {
				t.Fatalf("Unexpected no error: wantErr=%t", c.wantErr)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
