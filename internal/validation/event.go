package validation

import (
	"fmt"
	"time"

	"github.com/GitProger/go-telecom-2025/internal/model"
)

func checkType[T any](event *model.Event) error {
	if event.ExtraParams == nil {
		return fmt.Errorf("missing extra params for event %d", event.EventID)
	}
	if _, ok := event.ExtraParams.(T); !ok {
		return fmt.Errorf("invalid extra params type for event %d: %v", event.EventID, event.ExtraParams)
	}
	return nil
}

func Validate(event *model.Event) error {
	if event.EventType != model.IncomingEvent && event.EventType != model.OutgoingEvent {
		return fmt.Errorf("unknown event type: %d", event.EventType)
	}
	if event.CompetitorID < 1 {
		return fmt.Errorf("invalid competitor ID: %d", event.CompetitorID)
	}

	if (event.EventID < 1 || event.EventID > 11) && (event.EventID != model.EventDisqualified && event.EventID != model.EventFinished) {
		return fmt.Errorf("unknown event: %d", event.EventID)
	} else if event.EventID == model.EventOnRange || event.EventID == model.EventTargetHit {
		if err := checkType[int](event); err != nil {
			return err
		}

		if event.EventID == model.EventTargetHit {
			target := event.ExtraParams.(int)
			if target < 1 || target > 5 {
				return fmt.Errorf("biathlon target number must be from 1 to 5: %d", target)
			}
		}
	} else if event.EventID == model.EventStartTimeSet {
		if err := checkType[time.Time](event); err != nil {
			return err
		}
	} else if event.EventID == model.EventCannotContinue {
		if err := checkType[string](event); err != nil {
			return err
		}
	} else if event.ExtraParams != nil {
		return fmt.Errorf("unexpected extra params for event %d", event.EventID)
	}

	return nil
}
