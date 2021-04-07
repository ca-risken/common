package portscan

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/CyberAgent/mimosa-core/proto/finding"
)

func (n *NmapResult) GetFindings(projectID uint32, dataSource, data string) []*finding.FindingForUpsert {
	var ret []*finding.FindingForUpsert
	findingNmap := &finding.FindingForUpsert{
		Description:      n.GetDescription(),
		DataSource:       dataSource,
		DataSourceId:     n.GetDataSourceID(),
		ResourceName:     n.ResourceName,
		ProjectId:        projectID,
		OriginalScore:    n.GetScore(),
		OriginalMaxScore: 10.0,
		Data:             data,
	}
	ret = append(ret, findingNmap)
	for key, detail := range n.ScanDetail {
		if _, ok := httpCheckResult[key]; !ok {
			continue
		}
		if detail == true {
			addResult := httpCheckResult[key]
			ret = append(ret, &finding.FindingForUpsert{
				Description:      addResult.getDescription(n.Target, n.Port),
				DataSource:       dataSource,
				DataSourceId:     n.GetDataSourceID(),
				ResourceName:     n.ResourceName,
				ProjectId:        projectID,
				OriginalScore:    addResult.GetScore(),
				OriginalMaxScore: 10.0,
				Data:             data,
			})
		}
	}
	return ret
}

func (n *NmapResult) GetTags() []string {
	ret := []string{}
	if n.Service != "unknown" {
		ret = append(ret, n.Service)
	}
	/*
		for key, detail := range n.ScanDetail {
			if _, ok := httpCheckResult[key]; !ok {
				continue
			}
			if detail == true {
				addResult := httpCheckResult[key]
				ret = append(ret, addResult.Tag...)
			}
		}
	*/
	return ret
}

func (n *NmapResult) GetScore() float32 {
	status := n.Status
	protocol := n.Protocol
	port := n.Port
	if strings.ToUpper(status) == "CLOSED" || (strings.ToUpper(protocol) == "TCP" && strings.Index(strings.ToUpper(status), "FILTERED") > -1) {
		return 1.0
	}
	if strings.ToUpper(protocol) == "UDP" {
		return 6.0
	}
	switch port {
	case 22, 3306, 5432, 6379:
		return 8.0
	default:
		score := getScoreByScanDetail(n.ScanDetail)
		return score
	}
}

func (n *NmapResult) GetDescription() string {
	return fmt.Sprintf("%v is %v. protocol: %v, port %v", n.Target, n.Status, n.Protocol, n.Port)
}

func (n *NmapResult) GetDataSourceID() string {
	input := fmt.Sprintf("%v:%v:%v", n.Target, n.Protocol, n.Port)
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

func getScoreByScanDetail(detail map[string]interface{}) float32 {
	val, ok := detail["status"]
	if !ok {
		return 6.0
	}
	if strings.Index(val.(string), "401") > -1 || strings.Index(val.(string), "403") > -1 {
		return 1.0
	}
	return 6.0

}

// TODO 現状使用されていないTagの使用方法/要否の検討
type additionalCheckResult struct {
	Score       float32
	Tag         []string
	Description string
}

func (a additionalCheckResult) getDescription(target string, port int) string {
	ret := a.Description
	ret = strings.Replace(ret, "{TARGET}", target, 1)
	ret = strings.Replace(ret, "{PORT}", strconv.Itoa(port), 1)
	return ret
}

func (a additionalCheckResult) GetScore() float32 {
	return a.Score
}

var httpCheckResult = map[string]additionalCheckResult{
	"isHTTPOpenProxy":          additionalCheckResult{Score: 0.8, Tag: []string{"http"}, Description: "{TARGET} is Potentially OPEN proxy. port: {PORT}"},
	"isSSHEnabledPasswordAuth": additionalCheckResult{Score: 0.8, Tag: []string{"ssh"}, Description: "{TARGET} is supported password authentication. port: {PORT}"},
	"isSMTPOpenRelay":          additionalCheckResult{Score: 0.8, Tag: []string{"smtp"}, Description: "{TARGET} is an open relay. port: {PORT}"},
}
