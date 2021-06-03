package portscan

import (
	"reflect"
	"testing"
)

func TestGetFindings(t *testing.T) {
	cases := []struct {
		name             string
		projectID        uint32
		dataSource       string
		data             string
		nmapResult       *NmapResult
		numberOfFindings int
	}{
		{
			name:       "1 findings created",
			projectID:  1001,
			dataSource: "hogeDataSource",
			data:       "hogeData",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "open",
				Port:         8080,
			},
			numberOfFindings: 1,
		},
		{
			name:       "2 findings created",
			projectID:  1001,
			dataSource: "hogeDataSource",
			data:       "hogeData",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "open",
				Port:         8080,
				ScanDetail: map[string]interface{}{
					"isHTTPOpenProxy": true,
				},
			},
			numberOfFindings: 2,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			findings := c.nmapResult.GetFindings(c.projectID, c.dataSource, c.data)
			nof := len(findings)
			if c.numberOfFindings != nof {
				t.Fatalf("Unexpected number of findings: want=%v, got=%v", c.numberOfFindings, nof)
			}
		})
	}
}

func TestGetTags(t *testing.T) {
	cases := []struct {
		name       string
		nmapResult *NmapResult
		tags       []string
	}{
		{
			name: "No tags created",
			nmapResult: &NmapResult{
				Service: "unknown",
			},
			tags: []string{},
		},
		{
			name: "1 tag created",
			nmapResult: &NmapResult{
				Service: "http",
			},
			tags: []string{"http"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tags := c.nmapResult.GetTags()
			if !reflect.DeepEqual(c.tags, tags) {
				t.Fatalf("Unexpected tags: want=%v, got=%v", c.tags, tags)
			}
		})
	}
}

func TestGetScore(t *testing.T) {
	cases := []struct {
		name        string
		nmapResult  *NmapResult
		expectScore float32
	}{
		{
			name: "Score: 1.0 Status_Closed",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "closed",
				Port:         8080,
			},
			expectScore: 1.0,
		},
		{
			name: "Score: 6.0 UDP",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "udp",
				Status:       "open",
				Port:         8080,
			},
			expectScore: 6.0,
		},
		{
			name: "Score: 6.0 tcp/ssh port",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "open",
				Port:         22,
			},
			expectScore: 6.0,
		},
		{
			name: "Score: 8.0 tcp/dangerous port",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "open",
				Port:         3306,
			},
			expectScore: 8.0,
		},
		{
			name: "Score: 6.0 other",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "unknown",
				Status:       "open",
				Port:         37564,
			},
			expectScore: 6.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			score := c.nmapResult.GetScore()
			if c.expectScore != score {
				t.Fatalf("Unexpected score: want=%v, got=%v", c.expectScore, score)
			}
		})
	}
}

func TestGetScoreByScanDetail(t *testing.T) {
	cases := []struct {
		name        string
		detail      map[string]interface{}
		expectScore float32
	}{
		{
			name:        "Score: 6.0 Status doesn't exist.",
			detail:      map[string]interface{}{},
			expectScore: 6.0,
		},
		{
			name: "Score: 1.0 Status 401/403",
			detail: map[string]interface{}{
				"status": "403 Forbidden",
			},
			expectScore: 1.0,
		},
		{
			name: "Score: 6.0 Status without 401/403",
			detail: map[string]interface{}{
				"status": "200 OK",
			},
			expectScore: 6.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			score := getScoreByScanDetail(c.detail)
			if c.expectScore != score {
				t.Fatalf("Unexpected score: want=%v, got=%v", c.expectScore, score)
			}
		})
	}
}
