package portscan

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Ullaakut/nmap/v2"
)

type NmapResult struct {
	Port         int
	Protocol     string
	Target       string
	Status       string
	Service      string
	ResourceName string
	ScanDetail   map[string]interface{}
}

func Scan(target, protocol string, fPort, tPort int) ([]*NmapResult, error) {
	var ret []*NmapResult
	nmapResults, err := runNmap(target, protocol, fPort, tPort)
	if err != nil {
		return []*NmapResult{}, err
	}
	for _, result := range nmapResults {
		_ = result.analyzeResult()
		ret = append(ret, result)
	}
	return ret, nil
}

func runNmap(target, protocol string, fPort, tPort int) ([]*NmapResult, error) {
	var nmapResults []*NmapResult
	scanner, err := getScanner(target, protocol, fPort, tPort)
	if err != nil {
		return []*NmapResult{}, err
	}

	result, _, err := scanner.Run()
	if err != nil {
		return []*NmapResult{}, err
	}
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			nmapResults = append(nmapResults, &NmapResult{
				Port:     int(port.ID),
				Protocol: protocol,
				Target:   target,
				Status:   port.State.State,
				Service:  port.Service.Name,
			})
		}
	}
	return nmapResults, nil
}

func getScanner(host, protocol string, fPort, tPort int) (*nmap.Scanner, error) {
	if protocol == "tcp" {
		scanner, err := nmap.NewScanner(
			nmap.WithTargets(host),
			nmap.WithPorts(fmt.Sprintf("%v-%v", fPort, tPort)),
			nmap.WithServiceInfo(),
			nmap.WithSkipHostDiscovery(),
			//			nmap.WithSYNScan(),
			nmap.WithTimingTemplate(nmap.TimingAggressive),
		)
		if err != nil {
			return nil, err
		}
		return scanner, nil
	}
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(host),
		nmap.WithPorts(fmt.Sprintf("%v-%v", fPort, tPort)),
		nmap.WithServiceInfo(),
		nmap.WithSkipHostDiscovery(),
		//		nmap.WithUDPScan(),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
	)
	if err != nil {
		return nil, err
	}
	return scanner, nil
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
