package portscan

import (
	"fmt"
	"time"

	"github.com/Ullaakut/nmap/v2"
)

type NmapResult struct {
	Port         int                    `json:"port"`
	Protocol     string                 `json:"protocol"`
	Target       string                 `json:"target"`
	Status       string                 `json:"status"`
	Service      string                 `json:"service"`
	ResourceName string                 `json:"resource_name"`
	ScanDetail   map[string]interface{} `json:"scan_detail"`
}

// TODO これまで握りつぶしていたエラーを呼び出し元で確認するために追加したError。確認が終わり次第削除予定
type ResultAnalysisError struct {
	err error
}

func (e *ResultAnalysisError) Error() string {
	return fmt.Sprintf("failed to analyze portscan results, cause=%v", e.err)
}

func (e *ResultAnalysisError) Unwrap() error {
	return e.err
}

func wrapResultAnalysisError(err error) error {
	return &ResultAnalysisError{err}
}

func Scan(target, protocol string, fPort, tPort int) ([]*NmapResult, error) {
	var ret []*NmapResult
	nmapResults, err := runNmap(target, protocol, fPort, tPort)
	if err != nil {
		return []*NmapResult{}, err
	}
	fmt.Println("finish run nmap")
	time.Sleep(20 * time.Second)
	for _, result := range nmapResults {
		// TODO 握りつぶしていたエラーを呼び出し元で判定して確認できるようにするため、専用の型で返す。確認が終わり次第errをそのまま返すように変更予定
		err = result.analyzeResult()
		if err != nil {
			return nil, wrapResultAnalysisError(err)
		}
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
