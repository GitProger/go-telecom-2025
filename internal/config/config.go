package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Laps        int           `json:"laps"`        // Amount of laps for main distance
	LapLen      int           `json:"lapLen"`      // Length of each main lap
	PenaltyLen  int           `json:"penaltyLen"`  // Length of each penalty lap
	FiringLines int           `json:"firingLines"` // Number of firing lines per lap
	Start       time.Time     `json:"start"`       // Planned start time for the first competitor
	StartDelta  time.Duration `json:"startDelta"`  // Planned interval between starts
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path) // the config file is small, so `Unmarsal` instead of `Decoder`
	if err != nil {
		return nil, err
	}
	var config Config

	aux := &struct {
		Start      string `json:"start"`
		StartDelta string `json:"startDelta"`
		*Config
	}{
		Config: &config,
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}
	if config.Start, err = time.Parse(time.TimeOnly, aux.Start); err != nil {
		return nil, err
	}
	if t, err := time.Parse(time.TimeOnly, aux.StartDelta); err != nil {
		return nil, err
	} else {
		config.StartDelta = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second
	}

	return &config, nil
}
