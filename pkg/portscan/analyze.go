package portscan

func (n *NmapResult) analyzeResult() error {
	if n.Status == "closed" || n.Protocol != "tcp" {
		return nil
	}

	switch n.Service {
	case "ssh":
		data, err := analyzeSSH(n.Target, n.Port)
		if err != nil {
			return err
		}
		n.ScanDetail = data
		return nil
	case "smtp", "smtps", "submission":
		data, err := analyzeSMTP(n.Target, n.Port)
		if err != nil {
			return err
		}
		n.ScanDetail = data
		return nil
	}
	switch n.Protocol {
	case "tcp":
		switch n.Port {
		default:
			data, err := analyzeHTTP(n.Target, n.Port)
			if err != nil {
				return err
			}
			n.ScanDetail = data
		}
	case "udp":
		return nil
	}
	return nil
}
