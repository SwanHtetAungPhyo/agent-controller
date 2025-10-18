package utils

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/client"
)

func ParseHumanFriendlySchedule(schedule string) (*client.ScheduleCalendarSpec, error) {
	defaultSpec := &client.ScheduleCalendarSpec{
		Hour:   []client.ScheduleRange{{Start: 9}},
		Minute: []client.ScheduleRange{{Start: 0}},
	}

	switch schedule {
	case "daily-9am":
		return &client.ScheduleCalendarSpec{
			Hour:   []client.ScheduleRange{{Start: 9}},
			Minute: []client.ScheduleRange{{Start: 0}},
		}, nil
	case "daily-5pm":
		return &client.ScheduleCalendarSpec{
			Hour:   []client.ScheduleRange{{Start: 17}},
			Minute: []client.ScheduleRange{{Start: 0}},
		}, nil
	case "weekdays-9am":
		return &client.ScheduleCalendarSpec{
			Hour:      []client.ScheduleRange{{Start: 9}},
			Minute:    []client.ScheduleRange{{Start: 0}},
			DayOfWeek: []client.ScheduleRange{{Start: 1}, {Start: 2}, {Start: 3}, {Start: 4}, {Start: 5}}, // Mon-Fri
		}, nil
	case "weekdays-5pm":
		return &client.ScheduleCalendarSpec{
			Hour:      []client.ScheduleRange{{Start: 17}},
			Minute:    []client.ScheduleRange{{Start: 0}},
			DayOfWeek: []client.ScheduleRange{{Start: 1}, {Start: 2}, {Start: 3}, {Start: 4}, {Start: 5}}, // Mon-Fri
		}, nil
	case "hourly":
		return &client.ScheduleCalendarSpec{
			Minute: []client.ScheduleRange{{Start: 0}},
		}, nil
	case "every-30-minutes":
		return &client.ScheduleCalendarSpec{
			Minute: []client.ScheduleRange{{Start: 0}, {Start: 30}},
		}, nil
	case "every-15-minutes":
		return &client.ScheduleCalendarSpec{
			Minute: []client.ScheduleRange{{Start: 0}, {Start: 15}, {Start: 30}, {Start: 45}},
		}, nil
	case "weekly-monday-9am":
		return &client.ScheduleCalendarSpec{
			Hour:      []client.ScheduleRange{{Start: 9}},
			Minute:    []client.ScheduleRange{{Start: 0}},
			DayOfWeek: []client.ScheduleRange{{Start: 1}}, // Monday
		}, nil
	case "monthly-first-9am":
		return &client.ScheduleCalendarSpec{
			Hour:       []client.ScheduleRange{{Start: 9}},
			Minute:     []client.ScheduleRange{{Start: 0}},
			DayOfMonth: []client.ScheduleRange{{Start: 1}}, // 1st of month
		}, nil
	case "market-open": // 9:30 AM ET
		return &client.ScheduleCalendarSpec{
			Hour:      []client.ScheduleRange{{Start: 9}},
			Minute:    []client.ScheduleRange{{Start: 30}},
			DayOfWeek: []client.ScheduleRange{{Start: 1}, {Start: 2}, {Start: 3}, {Start: 4}, {Start: 5}}, // Mon-Fri
		}, nil
	case "market-close": // 4:00 PM ET
		return &client.ScheduleCalendarSpec{
			Hour:      []client.ScheduleRange{{Start: 16}},
			Minute:    []client.ScheduleRange{{Start: 0}},
			DayOfWeek: []client.ScheduleRange{{Start: 1}, {Start: 2}, {Start: 3}, {Start: 4}, {Start: 5}}, // Mon-Fri
		}, nil
	default:
		if parsedTime, err := time.Parse("15:04", schedule); err == nil {
			return &client.ScheduleCalendarSpec{
				Hour:   []client.ScheduleRange{{Start: parsedTime.Hour()}},
				Minute: []client.ScheduleRange{{Start: parsedTime.Minute()}},
			}, nil
		}
		return defaultSpec, fmt.Errorf("unknown schedule pattern: %s", schedule)
	}
}
