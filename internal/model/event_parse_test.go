package model_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/GitProger/go-telecom-2025/internal/model"
	"github.com/stretchr/testify/assert"
)

/**
Example input events from README.md:
[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000
[09:29:45.734] 3 1
[09:30:01.005] 4 1
[09:49:31.659] 5 1 1
[09:49:33.123] 6 1 1
[09:49:34.650] 6 1 2
[09:49:35.937] 6 1 4
[09:49:37.364] 6 1 5
[09:49:38.339] 7 1
[09:49:55.915] 8 1
[09:51:48.391] 9 1
[09:59:03.872] 10 1
[09:59:03.872] 11 1 Lost in the forest
*/

func tm(tm string) time.Time {
	t, err := time.Parse(model.TimeLayout, tm)
	if err != nil {
		panic(err)
	}
	return t
}

func TestParse(t *testing.T) {
	input := []string{
		"[09:05:59.867] 1 1",
		"[09:15:00.841] 2 1 09:30:00.000",
		"[09:29:45.734] 3 1",
		"[09:30:01.005] 4 1",
		"[09:49:31.659] 5 1 1",
		"[09:49:33.123] 6 1 1",
		"[09:49:34.650] 6 1 2",
		"[09:49:35.937] 6 1 4",
		"[09:49:37.364] 6 1 5",
		"[09:49:38.339] 7 1",
		"[09:49:55.915] 8 1",
		"[09:51:48.391] 9 1",
		"[09:59:03.872] 10 1",
		"[09:59:03.872] 11 1 Lost in the forest",
	}
	for i, test := range []struct {
		input      string
		output     *model.Event
		shouldFail bool
	}{
		{input: input[0], output: &model.Event{Time: tm("09:05:59.867"), EventID: 1, CompetitorID: 1}},
		{input: input[1], output: &model.Event{Time: tm("09:15:00.841"), EventID: 2, CompetitorID: 1, ExtraParams: tm("09:30:00.000")}},
		{input: input[2], output: &model.Event{Time: tm("09:29:45.734"), EventID: 3, CompetitorID: 1}},
		{input: input[3], output: &model.Event{Time: tm("09:30:01.005"), EventID: 4, CompetitorID: 1}},
		{input: input[4], output: &model.Event{Time: tm("09:49:31.659"), EventID: 5, CompetitorID: 1, ExtraParams: 1}},
		{input: input[5], output: &model.Event{Time: tm("09:49:33.123"), EventID: 6, CompetitorID: 1, ExtraParams: 1}},
		{input: input[6], output: &model.Event{Time: tm("09:49:34.650"), EventID: 6, CompetitorID: 1, ExtraParams: 2}},
		{input: input[7], output: &model.Event{Time: tm("09:49:35.937"), EventID: 6, CompetitorID: 1, ExtraParams: 4}},
		{input: input[8], output: &model.Event{Time: tm("09:49:37.364"), EventID: 6, CompetitorID: 1, ExtraParams: 5}},
		{input: input[9], output: &model.Event{Time: tm("09:49:38.339"), EventID: 7, CompetitorID: 1}},
		{input: input[10], output: &model.Event{Time: tm("09:49:55.915"), EventID: 8, CompetitorID: 1}},
		{input: input[11], output: &model.Event{Time: tm("09:51:48.391"), EventID: 9, CompetitorID: 1}},
		{input: input[12], output: &model.Event{Time: tm("09:59:03.872"), EventID: 10, CompetitorID: 1}},
		{input: input[13], output: &model.Event{Time: tm("09:59:03.872"), EventID: 11, CompetitorID: 1, ExtraParams: "Lost in the forest"}},

		{input: "[09:59:03.872] 100 1", shouldFail: true},
		{input: "[09:59:03.872] 100 abc", shouldFail: true},
		{input: "[09:59:03.872] 2 1 ", shouldFail: true},
		{input: "[09:59:03.872] 5 1", shouldFail: true},
		{input: "[09:59:03.872] 5 1 target", shouldFail: true},
		{input: "[09:59:03] 1 1", shouldFail: true},
		{input: "[xxx] 1 1", shouldFail: true},
	} {
		name := strconv.Itoa(i)
		if test.shouldFail {
			name += "_bad"
		}

		t.Run(name, func(t *testing.T) {
			event, err := model.ParseEvent(test.input)
			if test.shouldFail {
				assert.Error(t, err)
				assert.Nil(t, event)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.output, event)
			}
		})
	}
}
