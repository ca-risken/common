package xray

import (
	"github.com/aws/aws-xray-sdk-go/xray"
)

func InitXRay(config xray.Config) error {
	return xray.Configure(config)
}
