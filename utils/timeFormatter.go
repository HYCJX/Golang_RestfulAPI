package utils

import (
	"fmt"
	"log"
	"time"
)

var daysOfWeek = make(map[string]time.Weekday, 0)

func init() {
	// Initialise daysOfWeek map:
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := d.String()
		daysOfWeek[name] = d
		daysOfWeek[name[:3]] = d
	}
}

func StringToWeekday(weekdays []string) ([]time.Weekday, error) {
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

func FormatTimeString(timeString string) (string) {
	layout := "2006-01-02T15:04:05Z"
	dateTime, err := time.Parse(layout, timeString)
	if err != nil {
		log.Fatal(err)
	}
	return dateTime.Format("2006-01-02T15:04:05Z07:00")
}
