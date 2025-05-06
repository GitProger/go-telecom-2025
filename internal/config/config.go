package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Laps        int    `json:"laps"`        // Amount of laps for main distance
	LapLen      int    `json:"lapLen"`      // Length of each main lap
	PenaltyLen  int    `json:"penaltyLen"`  // Length of each penalty lap
	FiringLines int    `json:"firingLines"` // Number of firing lines per lap
	Start       string `json:"start"`       // Planned start time for the first competitor
	StartDelta  string `json:"startDelta"`  // Planned interval between starts
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path) // the config file is small, so `Unmarsal` instead of `Decoder`
	if err != nil {
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
