package config_test

import (
	"os"
	"testing"

	"github.com/GitProger/go-telecom-2025/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	const json = `{"laps" : 2,     "lapLen": 3651,   "penaltyLen": 50,"firingLines":1,"start":"09:30:00","startDelta": "00:00:30"}`

	tmpFile, err := os.CreateTemp("", "testfile-*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(json)
	assert.NoError(t, err)

	configOk := config.Config{
		Laps:        2,
		LapLen:      3651,
		PenaltyLen:  50,
		FiringLines: 1,
		Start:       "09:30:00",
		StartDelta:  "00:00:30",
	}

	configLoaded, err := config.LoadConfig(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, configOk, *configLoaded)
}
