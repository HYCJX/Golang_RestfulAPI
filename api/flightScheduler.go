package api

import (
	"log"
	"time"
)

type PilotInfo struct {
	Name     string
	Base     string
	Workdays []time.Weekday
	Flights  []Flight
}

type Flight struct {
	DepartDateTime time.Time
	ReturnDataTime time.Time
}

func newPilotInfo(name string, base string, workdays []string) *PilotInfo {
	timeWeekdays, err := stringToWeekday(workdays)
	if err != nil {
		log.Fatal(err)
	}
	pilotInfo := PilotInfo{Name: name, Base: base, Workdays: timeWeekdays, Flights: make([]Flight, 0)}
	return &pilotInfo
}
