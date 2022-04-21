package tracer

import (
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Config struct {
	ServiceName string
	Environment string
	Debug       bool
}

func Start(c *Config) {
	tracerOpts := []ddtracer.StartOption{
		ddtracer.WithEnv(c.Environment),
		ddtracer.WithService(c.ServiceName),
		ddtracer.WithDebugMode(c.Debug),
	}
	ddtracer.Start(tracerOpts...)
}

func Stop() {
	ddtracer.Stop()
}
