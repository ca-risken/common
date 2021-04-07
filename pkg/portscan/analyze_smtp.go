package portscan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Ullaakut/nmap"
)

func analyzeSMTP(target string, port int) (map[string]interface{}, error) {

	ret := map[string]interface{}{
		"isSMTPOpenRelay": checkSMTPOpenRelay(target, port),
	}
	return ret, nil
}

func checkSMTPOpenRelay(target string, port int) bool {
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
		return false
	}
	result, warn, err := scanner.Run()
	if err != nil {
		fmt.Printf("Nmap warning: %v", warn)
		return false
	}
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			for _, script := range port.Scripts {
				if strings.Index(script.Output, "Server is an open relay") > -1 {
					return true
				}
			}
		}
	}
	return false
}
