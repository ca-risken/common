package portscan

import (
	"fmt"

	"github.com/Ullaakut/nmap/v2"
)

type NmapResult struct {
	Port         int                    `json:"port"`
	Protocol     string                 `json:"protocol"`
	Target       string                 `json:"target"`
	Status       string                 `json:"status"`
	Service      string                 `json:"service"`
	ResourceName string                 `json:"resource_name"`
	ExternalLink string                 `json:"external_link"`
	ScanDetail   map[string]interface{} `json:"scan_detail"`
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

	result, warn, err := scanner.Run()
	if err != nil {
		fmt.Printf("Nmap warning: %v", warn)
		return []*NmapResult{}, err
	}
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			nmapResults = append(nmapResults, &NmapResult{
				Port:         int(port.ID),
				Protocol:     protocol,
				Target:       target,
				Status:       port.State.State,
				Service:      port.Service.Name,
				ExternalLink: makeURL(target, int(port.ID)),
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
			nmap.WithSYNScan(),
			nmap.WithTimingTemplate(nmap.TimingAggressive),
		)
		if err != nil {
			return nil, err
		}
		return scanner, nil
	}
	var scanner *nmap.Scanner
	var err error
	if fPort == 0 && tPort == 0 {
		scanner, err = nmap.NewScanner(
			nmap.WithTargets(host),
			nmap.WithServiceInfo(),
			nmap.WithSkipHostDiscovery(),
			nmap.WithUDPScan(),
			nmap.WithTimingTemplate(nmap.TimingAggressive),
		)
	} else {
		scanner, err = nmap.NewScanner(
			nmap.WithTargets(host),
			nmap.WithPorts(fmt.Sprintf("%v-%v", fPort, tPort)),
			nmap.WithServiceInfo(),
			nmap.WithSkipHostDiscovery(),
			nmap.WithUDPScan(),
			nmap.WithTimingTemplate(nmap.TimingAggressive),
		)
	}
	if err != nil {
		return nil, err
	}
	return scanner, nil
}
