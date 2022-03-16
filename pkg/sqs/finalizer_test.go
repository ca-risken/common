package sqs

import (
	"context"
	"errors"
	"testing"

	"github.com/ca-risken/core/proto/finding"
	"github.com/ca-risken/core/proto/finding/mocks"
	"github.com/stretchr/testify/mock"
)

var (
	sampleProjectID uint32 = 1
)

func TestFinal(t *testing.T) {
	type args struct {
		ProjectID  *uint32
		DataSource string
		Err        error
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
				ProjectID:  &sampleProjectID,
				DataSource: "google:scc",
				Err:        nil,
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
				ProjectID:  &sampleProjectID,
				DataSource: "google:scc",
				Err:        errors.New("Failed to scan"),
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
				ProjectID:  nil,
				DataSource: "google:scc",
				Err:        errors.New("Failed to scan"),
			},
			mockResp: mockResponse{},
			wantErr:  true,
		},
		{
			name: "ProjectID is nil(no error)",
			input: args{
				ProjectID:  nil,
				DataSource: "google:scc",
				Err:        nil,
			},
			mockResp: mockResponse{},
			wantErr:  false,
		},
		{
			name: "Failed to PutFinding API",
			input: args{
				ProjectID:  &sampleProjectID,
				DataSource: "google:scc",
				Err:        nil,
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
				ProjectID:  &sampleProjectID,
				DataSource: "google:scc",
				Err:        nil,
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
			f := Finalizer{findingClient: &mockClient}
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
			err := f.Final(ctx, c.input.ProjectID, c.input.DataSource, c.input.Err)
			if err == nil && c.wantErr {
				t.Fatalf("Unexpected no error: wantErr=%t", c.wantErr)
			}
			if err != nil && !c.wantErr {
				t.Fatalf("Unexpected error: wantErr=%t, err=%+v", c.wantErr, err)
			}
		})
	}
}
