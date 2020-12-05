package api

import (
	"fmt"
	"time"
)

var daysOfWeek = make(map[string]time.Weekday, 0)

func init() {
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := d.String()
		daysOfWeek[name] = d
		daysOfWeek[name[:3]] = d
	}
}

func stringToWeekday(weekdays []string) ([]time.Weekday, error) {
	var err error
	timeWeekdays := make([]time.Weekday, 0)
	for _, weekday := range weekdays {
		if timeWeekday, ok := daysOfWeek[weekday]; ok {
			timeWeekdays = append(timeWeekdays, timeWeekday)
		} else {
			err = fmt.Errorf("%w; invalid weekday '%s'", err, weekday)
		}
	}
	return timeWeekdays, err
}
