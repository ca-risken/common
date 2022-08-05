package portscan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Ullaakut/nmap"
)

func analyzeSMTP(target string, port int) (map[string]interface{}, error) {
	open, err := checkSMTPOpenRelay(target, port)
	if err != nil {
		return nil, err
	}
	ret := map[string]interface{}{
		"isSMTPOpenRelay": open,
	}
	return ret, nil
}

func checkSMTPOpenRelay(target string, port int) (bool, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithPorts(strconv.Itoa(port)),
		nmap.WithServiceInfo(),
		nmap.WithSkipHostDiscovery(),
		nmap.WithSYNScan(),
		nmap.WithScripts("smtp-open-relay"),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
	)
	if err != nil {
		return false, fmt.Errorf("failed to create scanner for SMTP, err=%w", err)
	}
	result, warn, err := scanner.Run()
	if err != nil {
		fmt.Printf("Nmap warning: %v", warn)
		return false, fmt.Errorf("failed to run scanner for SMTP, err=%w", err)
	}
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			for _, script := range port.Scripts {
				if strings.Contains(script.Output, "Server is an open relay") {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
