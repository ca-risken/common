package profiler

import (
	"fmt"

	ddprofiler "gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

type ProfileType int

const (
	ProfileTypeUndefined ProfileType = iota
	ProfileTypeCPUProfile
	ProfileTypeHeapProfile
	ProfileTypeBlockProfile
	ProfileTypeMutexProfile
	ProfileTypeGoroutineProfile
)

type ExporterType int

const (
	ExporterTypeUndefined ExporterType = iota
	ExporterTypeNOP
	ExporterTypeDatadog
)

type Config struct {
	ServiceName  string
	EnvName      string
	ProfileTypes []ProfileType
	ExporterType
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
			case ProfileTypeCPUProfile:
				pTypes = append(pTypes, ddprofiler.CPUProfile)
			case ProfileTypeHeapProfile:
				pTypes = append(pTypes, ddprofiler.HeapProfile)
			case ProfileTypeBlockProfile:
				pTypes = append(pTypes, ddprofiler.BlockProfile)
			case ProfileTypeMutexProfile:
				pTypes = append(pTypes, ddprofiler.MutexProfile)
			case ProfileTypeGoroutineProfile:
				pTypes = append(pTypes, ddprofiler.GoroutineProfile)
			default:
				return fmt.Errorf("undefined Profile Type: %v", c.ProfileTypes)
			}
		}
	}

	var err error
	switch c.ExporterType {
	case ExporterTypeDatadog:
		err = ddprofiler.Start(
			ddprofiler.WithService(c.ServiceName),
			ddprofiler.WithEnv(c.EnvName),
			ddprofiler.WithProfileTypes(pTypes...),
		)
	case ExporterTypeNOP:
		// nop
	default:
		err = fmt.Errorf("undefined Profile Exporter: %d", c.ExporterType)
	}
	return err
}

func (c *Config) Stop() {
	switch c.ExporterType {
	case ExporterTypeDatadog:
		ddprofiler.Stop()
	}
}

func ConvertProfileTypeFrom(typeStrings []string) ([]ProfileType, error) {
	var ret []ProfileType
	for _, ts := range typeStrings {
		switch ts {
		case "CPU":
			ret = append(ret, ProfileTypeCPUProfile)
		case "Heap":
			ret = append(ret, ProfileTypeHeapProfile)
		case "Block":
			ret = append(ret, ProfileTypeBlockProfile)
		case "Mutex":
			ret = append(ret, ProfileTypeMutexProfile)
		case "Goroutine":
			ret = append(ret, ProfileTypeGoroutineProfile)
		default:
			return ret, fmt.Errorf("undefined Profile Type: %s", ts)
		}
	}
	return ret, nil
}

func ConvertExporterTypeFrom(typeString string) (ExporterType, error) {
	switch typeString {
	case "nop":
		return ExporterTypeNOP, nil
	case "datadog":
		return ExporterTypeDatadog, nil
	default:
		return ExporterTypeUndefined, fmt.Errorf("undefined Profile Exporter Type: %s", typeString)
	}
}
