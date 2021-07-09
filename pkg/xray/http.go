package xray

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func IgnoreHealthCheckTracingHandler(segmentName, healthzPath string, next http.Handler) http.Handler {
	xrh := xray.Handler(xray.NewFixedSegmentNamer(segmentName), next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != healthzPath {
			xrh.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func AnnotateEnvTracingHandler(envName string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := xray.AddAnnotation(r.Context(), "env", envName); err != nil {
			// TODO logger
			// appLogger.Warnf("failed to annotate environment to x-ray: %+v", err)
			fmt.Printf("failed to annotate environment to x-ray: %+v", err)
		}
		next.ServeHTTP(w, r)
	})
}
