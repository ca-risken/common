package main

import (
	"fmt"

	"github.com/ca-risken/common/pkg/portscan"
)

func main() {
	for _, target := range scanTargets {
		results, err := portscan.Scan(target.target, target.protocol, target.fromPort, target.toPort)
		fmt.Printf("Scan start for: %v\n", target.target)
		if err != nil {
			fmt.Printf("err: %v\n", err)
		}
		for _, r := range results {
			fmt.Printf("nmapResult: %v\n", r)
			findings := r.GetFindings(1001, "dataSourceName", "dataBinary")
			for _, f := range findings {
				fmt.Printf("Finding: %v\n", f)
			}
			tags := r.GetTags()
			fmt.Printf("tags: %v\n", tags)
		}
	}
}

type target struct {
	target   string
	protocol string
	fromPort int
	toPort   int
}

var scanTargets = []target{
	{target: "google.com", protocol: "tcp", fromPort: 443, toPort: 443},
	// if you run open proxy in local machine, you can get OpenProxyEnabled finding.
	//	{target: "127.0.0.1", protocol: "tcp", fromPort: 8080, toPort: 8080},
	// if you run open relay smtp server in local machine, you can get OpenRelayEnabled finding.
	//	{target: "127.0.0.1", protocol: "tcp", fromPort: 25, toPort: 25},
	// if you run ssh allowed password authentication in local machine, you can get PasswordAuthEnabled finding.
	//	{target: "127.0.0.1", protocol: "tcp", fromPort: 22, toPort: 22},
}
