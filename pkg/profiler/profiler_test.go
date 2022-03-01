package profiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	// valid ProfileTypes
	c := &Config{
		ProfileTypes: []ProfileType{
			ProfileTypeCPUProfile,
			ProfileTypeHeapProfile,
			ProfileTypeBlockProfile,
			ProfileTypeMutexProfile,
			ProfileTypeGoroutineProfile,
		},
		ExporterType: ExporterTypeNOP,
	}
	err := c.Start()
	assert.NoError(t, err)

	// invalid ProfileTypes
	c = &Config{ProfileTypes: []ProfileType{
		ProfileType(8),
	}}
	err = c.Start()
	assert.Error(t, err)

	// invalid ExporterType
	c = &Config{ExporterType: 9}
	err = c.Start()
	assert.Error(t, err)
}

func TestConvertProfileTypeFrom(t *testing.T) {
	// invalid ProfileType string
	_, err := ConvertProfileTypeFrom([]string{"Undefined"})
	assert.Error(t, err)
}

func TestConvertExporterTypeFrom(t *testing.T) {
	// invalid ExporterType string
	_, err := ConvertExporterTypeFrom("Undefined")
	assert.Error(t, err)
}
