package profiler

import (
	"errors"
	"fmt"

	ddprofiler "gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type ProfileType int

const (
	Undefined ProfileType = iota
	CPUProfile
	HeapProfile
	BlockProfile
	MutexProfile
	GoroutineProfile
)

type Config struct {
	ServiceName  string
	EnvName      string
	ProfileTypes []ProfileType
	UseDatadog   bool
}

func (c *Config) Start() error {

	var pTypes []ddprofiler.ProfileType
	if len(c.ProfileTypes) == 0 {
		// default
		pTypes = []ddprofiler.ProfileType{
			ddprofiler.CPUProfile,
			ddprofiler.HeapProfile,
		}
	} else {
		for _, t := range c.ProfileTypes {
			switch t {
			case CPUProfile:
				pTypes = append(pTypes, ddprofiler.CPUProfile)
			case HeapProfile:
				pTypes = append(pTypes, ddprofiler.HeapProfile)
			case BlockProfile:
				pTypes = append(pTypes, ddprofiler.BlockProfile)
			case MutexProfile:
				pTypes = append(pTypes, ddprofiler.MutexProfile)
			case GoroutineProfile:
				pTypes = append(pTypes, ddprofiler.GoroutineProfile)
			default:
				return errors.New(fmt.Sprintf("undefined Profile Type: %v", c.ProfileTypes))
			}
		}
	}

	var err error
	if c.UseDatadog {
		err = ddprofiler.Start(
			ddprofiler.WithService(c.ServiceName),
			ddprofiler.WithEnv(c.EnvName),
			ddprofiler.WithProfileTypes(pTypes...),
		)
	}
	// nop when not specify datadog profiler
	return err
}

func (c *Config) Stop() {
	if c.UseDatadog {
		ddprofiler.Stop()
	}
}

func ConvertFrom(typeStrings []string) ([]ProfileType, error) {
	var ret []ProfileType
	for _, ts := range typeStrings {
		switch ts {
		case "CPU":
			ret = append(ret, CPUProfile)
		case "Heap":
			ret = append(ret, HeapProfile)
		case "Block":
			ret = append(ret, BlockProfile)
		case "Mutex":
			ret = append(ret, MutexProfile)
		case "Goroutine":
			ret = append(ret, GoroutineProfile)
		default:
			return ret, errors.New(fmt.Sprintf("undefined Profile Type: %s", ts))
		}
	}
	return ret, nil
}
