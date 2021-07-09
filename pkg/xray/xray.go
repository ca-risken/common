package xray

import (
	"github.com/aws/aws-xray-sdk-go/xray"
)

func InitXRay(config xray.Config) {
	xray.Configure(config)
}
