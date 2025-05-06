package model_test

import (
	"testing"
	"time"

	"github.com/GitProger/go-telecom-2025/internal/config"
	"github.com/GitProger/go-telecom-2025/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCompetitor(t *testing.T) {
	config := config.Config{
		Laps:        2,
		LapLen:      3651,
		PenaltyLen:  50,
		FiringLines: 1,
		Start:       "09:30:00",
		StartDelta:  "00:00:30",
	}
	comp := model.NewCompetitor(1, &config)
	comp.Status = model.NotFinished
	comp.Laps = []time.Duration{time.Duration((29*60 + 3.872) * float64(time.Second))}
	comp.PenaltyLaps = time.Duration(104.296 * float64(time.Second))
	comp.FiringLines = 1
	comp.Hits = 4
	assert.Equal(t, comp.String(), "[NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.479} 4/5")
	// 0.481 * 1:44.296 is 50.166376
	// while 104.296*0.480 = 50.06208
	//       104.296*0.479 = 49.95778
	// which is way closer, so there is an error in the README.md's example
}
