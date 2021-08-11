package portscan

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Ullaakut/nmap"
)

func analyzeHTTP(target string, port int) (map[string]interface{}, error) {
	url := makeURL(target, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Timeout:   time.Duration(5 * time.Second),
		Transport: transport,
	}
	var redirectURL []string
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		u := req.URL
		u.RawQuery = ""
		redirectURL = append(redirectURL, u.String())
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	resp, err := client.Do(req)
	if err != nil {
		// TODO add Logger
		fmt.Printf("%v", err)
		return nil, err
	}
	defer resp.Body.Close()

	ret := map[string]interface{}{
		"status":          resp.Status,
		"isHTTPOpenProxy": checkHTTPOpenProxy(target, port),
	}
	if resp.Header.Get("Server") != "" {
		ret["server"] = resp.Header.Get("Server")
	}
	if len(redirectURL) > 0 {
		ret["redirect"] = redirectURL
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
		// TODO add Logger
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

func makeURL(target string, port int) string {
	switch port {
	case 443:
		return fmt.Sprintf("https://%v", target)
	case 80:
		return fmt.Sprintf("http://%v", target)
	default:
		return fmt.Sprintf("http://%v:%v", target, port)
	}
}
