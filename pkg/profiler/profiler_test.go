package profiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	// valid ProfileTypes
	c := &Config{
		ProfileTypes: []ProfileType{
			CPUProfile,
			HeapProfile,
			BlockProfile,
			MutexProfile,
			GoroutineProfile,
		},
	}
	err := c.Start()
	assert.NoError(t, err)

	// invalid ProfileTypes
	c = &Config{ProfileTypes: []ProfileType{
		ProfileType(8),
	}}
	err = c.Start()
	assert.Error(t, err)
}

func TestConvertFrom(t *testing.T) {
	// invalid ProfileType string
	_, err := ConvertFrom([]string{"Undefined"})
	assert.Error(t, err)
}
