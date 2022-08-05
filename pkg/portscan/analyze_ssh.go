package portscan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Ullaakut/nmap"
)

func analyzeSSH(target string, port int) (map[string]interface{}, error) {
	open, err := checkSSHPasswordAuth(target, port)
	if err != nil {
		return nil, err
	}
	ret := map[string]interface{}{
		"isSSHEnabledPasswordAuth": open,
	}
	return ret, nil
}

func checkSSHPasswordAuth(target string, port int) (bool, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithPorts(strconv.Itoa(port)),
		nmap.WithServiceInfo(),
		nmap.WithSkipHostDiscovery(),
		nmap.WithSYNScan(),
		nmap.WithScripts("ssh-auth-methods"),
		nmap.WithTimingTemplate(nmap.TimingAggressive),
	)
	if err != nil {
		return false, fmt.Errorf("failed to create scanner for SSH, err=%w", err)
	}
	result, warn, err := scanner.Run()
	if err != nil {
		fmt.Printf("Nmap warning: %v", warn)
		return false, fmt.Errorf("failed to run scanner for SSH, err=%w", err)
	}
	for _, host := range result.Hosts {
		for _, port := range host.Ports {
			for _, script := range port.Scripts {
				if strings.Contains(script.Output, "password") {
					return true, nil
				}
			}
		}
	}
	return false, nil
}
