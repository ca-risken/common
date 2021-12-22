package portscan

import (
	"fmt"
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

func TestGetDescription(t *testing.T) {
	cases := []struct {
		name       string
		nmapResult *NmapResult
		expect     string
	}{
		{
			name: "Not Exist StatusDetail",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "closed",
				Service:      "hoge-service",
				Port:         8080,
			},
			expect: fmt.Sprintf("target: %v, protocol: %v, port: %v, status: %v, service: %v", "example.com", "tcp", "8080", "closed", "hoge-service"),
		},
		{
			name: "Exist StatusDetail (status)",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "closed",
				Service:      "hoge-service",
				Port:         8080,
				ScanDetail: map[string]interface{}{
					"status": 200,
				},
			},
			expect: fmt.Sprintf("target: %v, protocol: %v, port: %v, status: %v, service: %v, code: %v", "example.com", "tcp", "8080", "closed", "hoge-service", "200"),
		},
		{
			name: "Exist StatusDetail (status,server)",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "closed",
				Service:      "hoge-service",
				Port:         8080,
				ScanDetail: map[string]interface{}{
					"status": 200,
					"server": "hoge-server",
				},
			},
			expect: fmt.Sprintf("target: %v, protocol: %v, port: %v, status: %v, service: %v, code: %v, server: %v", "example.com", "tcp", "8080", "closed", "hoge-service", "200", "hoge-server"),
		},
		{
			name: "Exist StatusDetail (status,server,redirect)",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "closed",
				Service:      "hoge-service",
				Port:         8080,
				ScanDetail: map[string]interface{}{
					"status":   200,
					"server":   "hoge-server",
					"redirect": []string{"http://hoge1.com/fuga", "http://hoge2.com/fuga"},
				},
			},
			expect: fmt.Sprintf("target: %v, protocol: %v, port: %v, status: %v, service: %v, code: %v, server: %v, redirect: %v", "example.com", "tcp", "8080", "closed", "hoge-service", "200", "hoge-server", "http://hoge1.com/fuga,http://hoge2.com/fuga"),
		},
		{
			name: "Exist StatusDetail (status,server,redirect) Over 200 char",
			nmapResult: &NmapResult{
				ResourceName: "hogeResource",
				Target:       "example.com",
				Protocol:     "tcp",
				Status:       "closed",
				Service:      "hoge-service",
				Port:         8080,
				ScanDetail: map[string]interface{}{
					"status":   200,
					"server":   "hoge-server",
					"redirect": []string{"http://hoge1.com/fuga", "http://hoge2.com/fuga", "http://hoge2.com/fuga/piyopiyopiyo"},
				},
			},
			expect: fmt.Sprintf("target: %v, protocol: %v, port: %v, status: %v, service: %v, code: %v, server: %v, redirect: %v", "example.com", "tcp", "8080", "closed", "hoge-service", "200", "hoge-server", "http://hoge1.com/fuga,http://hoge2.com/fuga,http://hoge2.com/fuga/pi..."),
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			description := c.nmapResult.GetDescription()
			if c.expect != description {
				fmt.Printf("%v\n", len(description))
				t.Fatalf("Unexpected score: want=%v, got=%v", c.expect, description)
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
		service     string
		port        int
		detail      map[string]interface{}
		expectScore float32
	}{
		{
			name:        "Score: 6.0 Status doesn't exist.",
			service:     "http",
			port:        80,
			detail:      map[string]interface{}{},
			expectScore: 6.0,
		},
		{
			name:    "Score: 1.0 HTTPS Status 401/403",
			service: "https",
			port:    8443,
			detail: map[string]interface{}{
				"status": "403 Forbidden",
			},
			expectScore: 1.0,
		},
		{
			name:    "Score: 1.0 Port 443 HTTPS Status 401/403",
			service: "http-alt",
			port:    443,
			detail: map[string]interface{}{
				"status": "403 Forbidden",
			},
			expectScore: 1.0,
		},
		{
			name:    "Score: 6.0 HTTP Status 401/403",
			service: "http-alt",
			port:    8080,
			detail: map[string]interface{}{
				"status": "403 Forbidden",
			},
			expectScore: 6.0,
		},
		{
			name:    "Score: 6.0 Status without 401/403",
			service: "http-alt",
			port:    8080,
			detail: map[string]interface{}{
				"status": "200 OK",
			},
			expectScore: 6.0,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			score := getScoreByScanDetail(c.service, c.port, c.detail)
			if c.expectScore != score {
				t.Fatalf("Unexpected score: want=%v, got=%v", c.expectScore, score)
			}
		})
	}
}

func TestGetAdditionalCheck(t *testing.T) {
	cases := []struct {
		name          string
		key           string
		expect        AdditionalCheckResult
		expectIsFound bool
	}{
		{
			name: "AdditionalCheckResult is found.",
			key:  "isHTTPOpenProxy",
			expect: AdditionalCheckResult{Score: 0.8, Tag: []string{"http"},
				Description: "{TARGET} is Potentially OPEN proxy. port: {PORT}",
				Risk: `HTTP Open Proxies is Enabled.
	- Malicious client can use an open proxy to launch an attack that originates from the proxy server's IP.`,
				Recommendation: `Stop the open proxy.
	- Restrict target TCP and UDP port to trusted IP addresses.
	- Allow specific users to use the proxy by authenticating them.`},
			expectIsFound: true,
		},
		{
			name:          "AdditionalCheckResult is not found.",
			key:           "hogefugapiyo",
			expect:        AdditionalCheckResult{},
			expectIsFound: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			acr, isFound := GetAdditionalCheckResult(c.key)
			if c.expectIsFound != isFound {
				t.Fatalf("Unexpected isFound: want=%v, got=%v", c.expectIsFound, isFound)
			}
			if !reflect.DeepEqual(c.expect, acr) {
				t.Fatalf("Unexpected AdditionalCheckResult: want=%v, got=%v", c.expect, acr)
			}
		})
	}
}

func TestGetAdditionalDescription(t *testing.T) {
	cases := []struct {
		name   string
		adr    AdditionalCheckResult
		target string
		port   int
		expect string
	}{
		{
			name:   "Get Description",
			adr:    AdditionalCheckResult{Score: 0.8, Tag: []string{"http"}, Description: "{TARGET} is Potentially OPEN proxy. port: {PORT}"},
			target: "hoge",
			port:   80,
			expect: "hoge is Potentially OPEN proxy. port: 80",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			desc := c.adr.GetDescription(c.target, c.port)
			if c.expect != desc {
				t.Fatalf("Unexpected description: want=%v, got=%v", c.expect, desc)
			}
		})
	}
}

func TestGetAdditionalScore(t *testing.T) {
	cases := []struct {
		name   string
		adr    AdditionalCheckResult
		expect float32
	}{
		{
			name:   "Get Description",
			adr:    AdditionalCheckResult{Score: 0.8, Tag: []string{"http"}, Description: "{TARGET} is Potentially OPEN proxy. port: {PORT}"},
			expect: 0.8,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			score := c.adr.GetScore()
			if c.expect != score {
				t.Fatalf("Unexpected score: want=%v, got=%v", c.expect, score)
			}
		})
	}
}

func TestGetAdditionaRisk(t *testing.T) {
	cases := []struct {
		name   string
		adr    AdditionalCheckResult
		expect string
	}{
		{
			name: "Get Risk",
			adr: AdditionalCheckResult{Score: 0.8, Tag: []string{"http"},
				Description: "{TARGET} is Potentially OPEN proxy. port: {PORT}",
				Risk: `HTTP Open Proxies is Enabled.
	- Malicious client can use an open proxy to launch an attack that originates from the proxy server's IP.`,
				Recommendation: `Stop the open proxy.
	- Restrict target TCP and UDP port to trusted IP addresses.
	- Allow specific users to use the proxy by authenticating them.`},
			expect: `HTTP Open Proxies is Enabled.
	- Malicious client can use an open proxy to launch an attack that originates from the proxy server's IP.`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			risk := c.adr.GetRisk()
			if c.expect != risk {
				t.Fatalf("Unexpected risk: want=%v, got=%v", c.expect, risk)
			}
		})
	}
}

func TestGetAdditionaRecommendation(t *testing.T) {
	cases := []struct {
		name   string
		adr    AdditionalCheckResult
		expect string
	}{
		{
			name: "Get Recommendation",
			adr: AdditionalCheckResult{Score: 0.8, Tag: []string{"http"},
				Description: "{TARGET} is Potentially OPEN proxy. port: {PORT}",
				Risk: `HTTP Open Proxies is Enabled.
	- Malicious client can use an open proxy to launch an attack that originates from the proxy server's IP.`,
				Recommendation: `Stop the open proxy.
	- Restrict target TCP and UDP port to trusted IP addresses.
	- Allow specific users to use the proxy by authenticating them.`},
			expect: `Stop the open proxy.
	- Restrict target TCP and UDP port to trusted IP addresses.
	- Allow specific users to use the proxy by authenticating them.`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			risk := c.adr.GetRecommendation()
			if c.expect != risk {
				t.Fatalf("Unexpected recommendation: want=%v, got=%v", c.expect, risk)
			}
		})
	}
}
