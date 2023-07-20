package tracer

import (
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Config struct {
	ServiceName  string
	Environment  string
	Debug        bool
	SamplingRate *float64 // Set a sampling rate value in the range of 0.0000 to 1.0000. (e.g.) 0.3 = 30%
}

func Start(c *Config) {
	rate := 1.0000 // 100%
	if c.SamplingRate != nil {
		rate = *c.SamplingRate
	}
	tracerOpts := []ddtracer.StartOption{
		ddtracer.WithEnv(c.Environment),
		ddtracer.WithService(c.ServiceName),
		ddtracer.WithDebugMode(c.Debug),
		ddtracer.WithSamplingRules(
			[]ddtracer.SamplingRule{
				ddtracer.ServiceRule(c.ServiceName, rate),
			},
		),
	}
	ddtracer.Start(tracerOpts...)
}

func Stop() {
	ddtracer.Stop()
}
