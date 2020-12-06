package utils

import (
	"math"
	"time"
)

var weekdayBitmap = map[time.Weekday]int{
	time.Sunday:    1,
	time.Monday:    2,
	time.Tuesday:   3,
	time.Wednesday: 4,
	time.Thursday:  5,
	time.Friday:    6,
	time.Saturday:  7,
}
var weekdayList = []time.Weekday{
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
	time.Sunday,
}

func EncodeWeekdays(weekdays []time.Weekday) int {
	encodedWeekdays := 0
	for _, weekday := range weekdays {
		encodedWeekdays += int(math.Pow(2, float64(weekdayBitmap[weekday])))
	}
	return encodedWeekdays
}

func DecodeWeekdays(encodedWeekdays int) []string {
	weekdays := make([]string, 0)
	for _, weekday := range weekdayList {
		if encodedWeekdays&int(math.Pow(2, float64(weekdayBitmap[weekday]))) != 0 {
			weekdays = append(weekdays, weekday.String())
		}
	}
	return weekdays
}
