package provider

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/GitProger/go-telecom-2025/internal/model"
)

func ScanFile(filename string) ([]*model.Event, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []*model.Event
	var last *model.Event

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()

		cur, err := model.ParseEvent(line)
		if err != nil {
			return nil, fmt.Errorf("parsing error: %w", err)
		}

		if last != nil && last.Time.After(cur.Time) {
			return nil, fmt.Errorf("event order error: %s > %s",
				last.Time.Format(model.TimeLayout),
				cur.Time.Format(model.TimeLayout))
		}

		last = cur
		events = append(events, cur)
	}

	return events, scanner.Err()
}

func Scan(ctx context.Context, source io.Reader) (<-chan *model.Event, <-chan error) {
	events := make(chan *model.Event)
	errs := make(chan error, 1)

	go func() {
		defer close(events)
		defer close(errs)
		var last *model.Event
		scanner := bufio.NewScanner(source)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errs <- ctx.Err()
				return
			default:
				line := scanner.Text()
				if line == "" {
					continue
				}

				cur, err := model.ParseEvent(line)
				if err != nil {
					errs <- fmt.Errorf("parsing error: %w", err)
					return
				}
				if last != nil && last.Time.After(cur.Time) {
					errs <- fmt.Errorf("event order error: %s > %s",
						last.Time.Format(model.TimeLayout),
						cur.Time.Format(model.TimeLayout))
					return
				}
				last = cur

				select {
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				case events <- cur:
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("scanner error: %w", err)
		}
	}()
	return events, errs
}
