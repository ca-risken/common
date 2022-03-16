package sqs

import "fmt"

const (
	DataSourceAWSGuardDuty      = "aws:guard-duty"
	DataSourceAWSAccessAnalyzer = "aws:access-analyzer"
	DataSourceAWSAdminChecker   = "aws:admin-checker"
	DataSourceAWSCloudsploit    = "aws:cloudsploit"
	DataSourceAWSPortscan       = "aws:portscan"
	DataSourceAWSActivity       = "aws:activity"

	DataSourceGoogleAsset       = "google:asset"
	DataSourceGoogleCloudSploit = "google:cloudsploit"
	DataSourceGoogleSCC         = "google:scc"
	DataSourceGooglePortscan    = "google:portscan"

	DataSourceCodeGitleaks = "code:gitleaks"

	DataSourceOSINTSubdomain = "osint:subdomain"
	DataSourceOSINTWebsite   = "osint:website"

	DataSourceDiagnosisWpscan   = "diagnosis:wpscan"
	DataSourceDiagnosisPortscan = "diagnosis:portscan"
	DataSourceDiagnosisAppscan  = "diagnosis:application-scan"
)

type recommend struct {
	Risk           string `json:"risk,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

func getRecommend(recommendType string) recommend {
	return recommendMap[recommendType]
}

var recommendMap = map[string]recommend{
	DataSourceAWSGuardDuty: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceAWSGuardDuty),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/aws/overview_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceAWSAccessAnalyzer: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceAWSAccessAnalyzer),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/aws/overview_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceAWSAdminChecker: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceAWSAdminChecker),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/aws/overview_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceAWSCloudsploit: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceAWSCloudsploit),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/aws/overview_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceAWSPortscan: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceAWSPortscan),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/aws/overview_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},

	DataSourceGoogleAsset: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceGoogleAsset),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/google/overview_gcp/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceGoogleCloudSploit: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceGoogleCloudSploit),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/google/overview_gcp/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceGoogleSCC: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceGoogleSCC),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/google/overview_gcp/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceGooglePortscan: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceGooglePortscan),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/google/overview_gcp/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},

	DataSourceCodeGitleaks: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceCodeGitleaks),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the access rights you set for the DataSource and the reachability of the network.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/code/gitleaks_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},

	DataSourceOSINTSubdomain: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceOSINTSubdomain),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/osint/datasource/
		- For Domain type, make sure the FQDN format is registered.
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceOSINTWebsite: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceOSINTWebsite),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the network is reachable to the target host.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/osint/datasource/
		- For Website type, make sure the URL format(e.g. http://example.com ) is registerd.
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},

	DataSourceDiagnosisWpscan: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceDiagnosisWpscan),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the network is reachable to the target host.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/code/gitleaks_datasource/
		- And please also check the FAQ page.
		- https://docs.security-hub.jp/contact/faq/#wpscan
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceDiagnosisPortscan: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceDiagnosisPortscan),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the network is reachable to the target host.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/code/gitleaks_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
	DataSourceDiagnosisAppscan: {
		Risk: fmt.Sprintf("Failed to scan %s, So you are not gathering the latest security threat information.", DataSourceDiagnosisAppscan),
		Recommendation: `Please review the following items and rescan,
		- Ensure the error message of the DataSource.
		- Ensure the network is reachable to the target host.
		- Refer to the documentation to make sure you have not omitted any of the steps you have set up.
		- https://docs.security-hub.jp/code/gitleaks_datasource/
		- If this does not resolve the problem, or if you suspect that the problem is server-side, please contact the system administrators.`,
	},
}
