package model

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const TimeLayout = "15:04:05.000" // like StampMilli in time package

const (
	EventRegister       = 1  // The competitor registered
	EventStartTimeSet   = 2  // The start time was set by a draw {startTime}
	EventOnStartLine    = 3  // The competitor is on the start line
	EventStarted        = 4  // The competitor has started
	EventOnRange        = 5  // The competitor is on the firing range {firingRange}
	EventTargetHit      = 6  // The target has been hit {target}
	EventLeftRange      = 7  // The competitor left the firing range
	EventEnteredPenalty = 8  // The competitor entered the penalty laps
	EventLeftPenalty    = 9  // The competitor left the penalty laps
	EventLapCompleted   = 10 // The competitor ended the main lap
	EventCannotContinue = 11 // The competitor can`t continue {comment}

	EventDisqualified = 32 // The competitor is disqualified
	EventFinished     = 33 // The competitor has finished
)

/*

Example:
[09:05:59.867] The competitor(1) registered
[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000
[09:29:45.734] The competitor(1) is on the start line
[09:30:01.005] The competitor(1) has started
[09:49:31.659] The competitor(1) is on the firing range(1)
[09:49:33.123] The target(1) has been hit by competitor(1)
[09:49:34.650] The target(2) has been hit by competitor(1)
[09:49:35.937] The target(4) has been hit by competitor(1)
[09:49:37.364] The target(5) has been hit by competitor(1)
[09:49:38.339] The competitor(1) left the firing range
[09:49:55.915] The competitor(1) entered the penalty laps
[09:51:48.391] The competitor(1) left the penalty laps
[09:59:03.872] The competitor(1) ended the main lap
[09:59:05.321] The competitor(1) can`t continue: Lost in the forest
*/

var eventComms = map[int]string{
	EventRegister:       "The competitor(%d) registered",
	EventStartTimeSet:   "The start time for the competitor(%d) was set by a draw to %s",
	EventOnStartLine:    "The competitor(%d) is on the start line",
	EventStarted:        "The competitor(%d) has started",
	EventOnRange:        "The competitor(%d) is on the firing range(%d)",
	EventTargetHit:      "The target(%d) has been hit by competitor(%d)",
	EventLeftRange:      "The competitor(%d) left the firing range",
	EventEnteredPenalty: "The competitor(%d) entered the penalty laps",
	EventLeftPenalty:    "The competitor(%d) left the penalty laps",
	EventLapCompleted:   "The competitor(%d) ended the main lap",
	EventCannotContinue: "The competitor(%d) can`t continue: %s",

	EventDisqualified: "The competitor(%d) is disqualified",
	EventFinished:     "The competitor(%d) has finished",
}

type EventType int

const (
	IncomingEvent EventType = iota
	OutgoingEvent
)

type Event struct {
	EventType    EventType // in/out
	EventID      int       // exact type of event
	CompetitorID int
	Time         time.Time
	ExtraParams  any
}

func (e *Event) String() string {
	format := eventComms[e.EventID]
	var outer string
	switch e.EventID {
	case EventStartTimeSet: // competitor number, start time
		outer = fmt.Sprintf(format, e.CompetitorID, e.ExtraParams.(time.Time).Format(TimeLayout))
	case EventOnRange: // competitor number, firing range number
		outer = fmt.Sprintf(format, e.CompetitorID, e.ExtraParams.(int))
	case EventTargetHit: // target number, competitor number
		outer = fmt.Sprintf(format, e.ExtraParams.(int), e.CompetitorID)
	case EventCannotContinue: // competitor numner, comment
		outer = fmt.Sprintf(format, e.CompetitorID, e.ExtraParams.(string))
	default: // just competitor number
		outer = fmt.Sprintf(format, e.CompetitorID)
	}
	return fmt.Sprintf("[%s] %s", e.Time.Format(TimeLayout), outer)
}

func ParseEvent(line string) (*Event, error) { // Incoming event only
	var event Event
	event.EventType = IncomingEvent
	var tm, extra string

	n, err := fmt.Sscanf(line, "%s %d %d %s", &tm, &event.EventID, &event.CompetitorID, &extra)
	if err != nil && (err != io.EOF && n < 4) {
		return nil, err
	}
	if event.Time, err = time.Parse("["+TimeLayout+"]", tm); err != nil {
		return nil, err
	}

	switch event.EventID {
	case EventStartTimeSet: // start time
		if event.ExtraParams, err = time.Parse(TimeLayout, extra); err != nil {
			return nil, err
		}
	case EventOnRange, EventTargetHit: // firing range number || target number
		if event.ExtraParams, err = strconv.Atoi(extra); err != nil {
			return nil, err
		}
	case EventCannotContinue: // comment
		parts := strings.SplitN(line, " ", 4) // [time] eventID competitorID comment
		event.ExtraParams = parts[3]
	case EventRegister, EventOnStartLine, EventStarted, EventLeftRange, EventEnteredPenalty, EventLeftPenalty, EventLapCompleted:
	case EventDisqualified, EventFinished: // outgoing event
		return nil, fmt.Errorf("outgoing event %d can not be parsed", event.EventID)
	default: // unknown event
		return nil, fmt.Errorf("unknown event type: %d", event.EventID)
	}
	return &event, nil
}
