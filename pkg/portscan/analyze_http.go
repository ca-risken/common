package portscan

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ullaakut/nmap"
)

func analyzeHTTP(target string, port int) (map[string]interface{}, error) {
	var url string
	if port == 443 {
		url = fmt.Sprintf("https://%v", target)
	} else if port == 80 {
		url = fmt.Sprintf("http://%v", target)
	} else {
		url = fmt.Sprintf("http://%v:%v", target, port)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   time.Duration(5 * time.Second),
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	ret := map[string]interface{}{
		"status":          resp.Status,
		"header":          resp.Header,
		"isHTTPOpenProxy": checkHTTPOpenProxy(target, port),
	}
	return ret, nil
}

func checkHTTPOpenProxy(target string, port int) bool {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(target),
		nmap.WithPorts(strconv.Itoa(port)),
		nmap.WithServiceInfo(),
		nmap.WithSkipHostDiscovery(),
		nmap.WithSYNScan(),
		nmap.WithScripts("http-open-proxy"),
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
				if strings.Index(script.Output, "Potentially OPEN proxy.") > -1 {
					return true
				}
			}
		}
	}
	return false
}
