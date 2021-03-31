package portscan

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func (n *NmapResult) analyzeResult() error {
	if n.Status == "closed" {
		return nil
	}
	switch n.Port {
	case 22:
		return nil
	default:
		data, err := analyzeHTTP(n.Target, n.Port)
		if err != nil {
			return nil
		}
		n.ScanDetail = data
	}
	return nil
}

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
		"status": resp.Status,
		"header": resp.Header,
	}
	return ret, nil
}
