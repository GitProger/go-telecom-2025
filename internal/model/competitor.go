package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/GitProger/go-telecom-2025/internal/config"
)

type CompetitorStatus int

const (
	NotStarted CompetitorStatus = iota
	Started
	NotFinished
	Finished
)

func formatDuration(d time.Duration) string {
	var z time.Time
	return z.Add(d).Format(TimeLayout)
}

type Competitor struct {
	config       *config.Config
	Arrived      bool
	IsFiring     bool
	Disqualified bool

	ID int

	PenaltyStartTime time.Time
	LapStartTime     time.Time

	PlannedStartTime time.Time
	StartTime        time.Time

	Status      CompetitorStatus
	Hits        int // successful hits
	FiringLines int // total: 5 * lines shots
	Laps        []time.Duration
	PenaltyLaps time.Duration // considered as one lap
	// totally penalty laps: number of misses = 5 * firing lines - hits
}

func NewCompetitor(id int, conf *config.Config) *Competitor {
	return &Competitor{
		config: conf,
		ID:     id,
		Status: NotStarted,
	}
}

func penaltyRange(comp *Competitor) int {
	return comp.config.PenaltyLen * (comp.FiringLines*5 - comp.Hits)
}

// The final report for each competitor:
// - Total time includes the difference between scheduled and actual start time or **NotStarted**/**NotFinished** marks
// - Time taken to complete each lap
// - Average speed for each lap [m/s]
// - Time taken to complete penalty laps
// - Average speed over penalty laps [m/s]
// - Number of hits/number of shots
// return example: [NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5
func (c *Competitor) String() string {
	var status string
	if st := c.Status; st == Finished {
		status = formatDuration(c.TimeFromPlannedStart())
	} else if st == NotStarted {
		status = "NotStarted"
	} else if st == NotFinished {
		status = "NotFinished"
	} else if st == Started {
		panic("can not call .String() on running competitor")
	}

	lapStr := func(length int, lapTime time.Duration) string {
		if lapTime != 0 {
			sp := float64(length) / lapTime.Seconds()
			return fmt.Sprintf("{%s, %.3f}", formatDuration(lapTime), math.Floor(sp*1000)/float64(1000))
		} else {
			return "{,}"
		}
	}

	var sb strings.Builder
	sb.WriteByte('[')
	for i := 0; i < c.config.Laps; i++ {
		if i < len(c.Laps) {
			sb.WriteString(lapStr(c.config.LapLen, c.Laps[i]))
		} else {
			sb.WriteString("{,}")
		}
		if i < c.config.Laps-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteByte(']')

	return fmt.Sprintf("[%s] %d %s %s %d/%d",
		status,
		c.ID,
		sb.String(),
		lapStr(penaltyRange(c), c.PenaltyLaps),
		c.Hits,
		c.FiringLines*5)
}

func (c *Competitor) TimeFromPlannedStart() time.Duration {
	return c.LapStartTime.Sub(c.PlannedStartTime)
}
