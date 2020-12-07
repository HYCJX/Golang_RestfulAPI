package utils

import (
	"math"
	"time"
)

var weekdayBitmap = map[time.Weekday]int{
	time.Sunday:    0,
	time.Monday:    1,
	time.Tuesday:   2,
	time.Wednesday: 3,
	time.Thursday:  4,
	time.Friday:    5,
	time.Saturday:  6,
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

// Encode a slice of time.Weekday.
// e.g. [Monday, Sunday] is encoded as binary [0,0,0,0,0,1,1], which is 3.
func EncodeWeekdays(weekdays []time.Weekday) int {
	encodedWeekdays := 0
	for _, weekday := range weekdays {
		encodedWeekdays += int(math.Pow(2, float64(weekdayBitmap[weekday])))
	}
	return encodedWeekdays
}

// Encode one string representation of date-time
func EncodeWeekday(dateTimeString string) (int, error) {
	layout := "2006-01-02T15:04:05Z"
	dateTime, err := time.Parse(layout, dateTimeString)
	if err != nil {
		return 0, err
	}
	return int(math.Pow(2, float64(weekdayBitmap[dateTime.Weekday()]))), err
}

// Decode from encoded int and convert to a slice of string representations of weekdays.
func DecodeWeekdays(encodedWeekdays int) []string {
	weekdays := make([]string, 0)
	for _, weekday := range weekdayList {
		if encodedWeekdays&int(math.Pow(2, float64(weekdayBitmap[weekday]))) != 0 {
			weekdays = append(weekdays, weekday.String())
		}
	}
	return weekdays
}
