package monitor

import (
	"fmt"
	"time"

	"github.com/GitProger/go-telecom-2025/internal/config"
	"github.com/GitProger/go-telecom-2025/internal/model"
	"github.com/GitProger/go-telecom-2025/internal/service"
)

type EventMonitor interface {
	DigestEvent(event *model.Event) (*model.Event, error)
	GetReport() []*model.Competitor
	Disqualified() []*model.Event
}

type monitor struct {
	lastTime     time.Time
	disqualified []int
	delta        time.Duration
	defaultStart time.Time

	conf    *config.Config
	service *service.CompetitorService
}

func NewEventMonitor(conf *config.Config) *monitor {
	return &monitor{
		conf:    conf,
		service: service.NewCompetitorService(),
	}
}

func (em *monitor) DigestEvent(event *model.Event) (*model.Event, error) {
	if em.delta == 0 {
		t, err := time.Parse("15:04:05", em.conf.StartDelta)
		if err != nil {
			return nil, fmt.Errorf("delta parsing error in config: %w", err)
		}
		em.delta = time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second

		em.defaultStart, err = time.Parse(model.TimeLayout, em.conf.Start)
		if err != nil {
			return nil, fmt.Errorf("start time parsing error in config: %w", err)
		}
	}

	em.lastTime = event.Time
	cId := event.CompetitorID
	comp := em.service.Get(cId)
	switch event.EventID {
	case model.EventRegister:
		em.service.Register(cId, em.conf)
	case model.EventStartTimeSet:
		comp.PlannedStartTime = event.ExtraParams.(time.Time)
	case model.EventOnStartLine:
		comp.Arrived = true
	case model.EventStarted:
		if !comp.Arrived {
			return nil, fmt.Errorf("competitor %d is not on start line", cId)
		}
		if comp.Disqualified {
			return nil, nil
		}
		comp.Status = model.Started
		comp.StartTime = event.Time
		comp.LapStartTime = comp.PlannedStartTime
	case model.EventOnRange:
		if event.ExtraParams.(int) != comp.FiringLines+1 {
			return nil, fmt.Errorf("competitor %d is on range %d, not %d", cId, event.ExtraParams.(int), comp.FiringLines+1)
		}
		comp.IsFiring = true
	case model.EventTargetHit:
		if !comp.IsFiring {
			return nil, fmt.Errorf("competitor %d is not firing", cId)
		}
		comp.Hits += 1
	case model.EventLeftRange:
		comp.FiringLines += 1
		comp.IsFiring = false
	case model.EventEnteredPenalty:
		comp.PenaltyStartTime = event.Time
	case model.EventLeftPenalty:
		if comp.PenaltyStartTime.IsZero() {
			return nil, fmt.Errorf("competitor %d left penalty area without entering it", cId)
		}
		comp.PenaltyLaps += event.Time.Sub(comp.PenaltyStartTime)
		comp.PenaltyStartTime = time.Time{}

	case model.EventLapCompleted: // includes penalty laps and shooting
		if len(comp.Laps) == em.conf.Laps {
			return nil, fmt.Errorf("competitor %d has already finished", cId)
		}

		lapTime := event.Time.Sub(comp.LapStartTime)
		comp.LapStartTime = event.Time // Finish time
		comp.Laps = append(comp.Laps, lapTime)

		if len(comp.Laps) == em.conf.Laps {
			comp.Status = model.Finished
			return &model.Event{
				EventType:    model.OutgoingEvent,
				EventID:      model.EventFinished,
				CompetitorID: cId,
				Time:         event.Time,
			}, nil
		}
	case model.EventCannotContinue:
		comp.Status = model.NotFinished
	}

	if id := em.findLate(); id != 0 {
		return em.disqualify(id), nil
	}

	return nil, nil
}

func (em *monitor) GetReport() []*model.Competitor {
	return em.service.GetAll()
}

func (em *monitor) Disqualified() []*model.Event {
	events := make([]*model.Event, len(em.disqualified))
	for i, id := range em.disqualified {
		events[i] = em.disqualify(id)
	}
	return events
}

func (em *monitor) disqualify(id int) *model.Event {
	return &model.Event{
		EventType:    model.OutgoingEvent,
		EventID:      model.EventDisqualified,
		CompetitorID: id,
		Time:         em.lastTime,
	}
}

func (em *monitor) findLate() int {
	for id, comp := range em.service.GetAllMap() {
		if comp.Disqualified {
			continue
		}
		if comp.Status == model.NotStarted {
			if em.lastTime.Sub(comp.PlannedStartTime) > em.delta {
				comp.Disqualified = true
				em.disqualified = append(em.disqualified, id)
			}
		}
	}

	if len(em.disqualified) > 0 {
		n := len(em.disqualified) - 1
		last := em.disqualified[n]
		em.disqualified = em.disqualified[:n]
		return last
	}
	return 0
}
