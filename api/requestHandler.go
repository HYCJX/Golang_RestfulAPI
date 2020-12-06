package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

var pilotInfoMap = make(map[int]*PilotInfo, 0)

func init() {
	// Initialise time-string converter:
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := d.String()
		daysOfWeek[name] = d
		daysOfWeek[name[:3]] = d
	}
	// Read json file:
	crew, err := readJson()
	if err != nil {
		log.Fatal(err)
	}
	// Fill in pilots-info:
	for _, pilot := range crew.Crew {
		pilotInfoMap[pilot.ID] = NewPilotInfo(pilot.ID, pilot.Name, pilot.Base, pilot.Workdays)
	}
}

type PilotInfo struct {
	ID       int         `json:"ID"`
	Name     string         `json:"Name"`
	Base     string         `json:"Base"`
	Workdays []time.Weekday `json:"Workdays"`
	Flights  []Flight       `json:"Flights"`
}

type Flight struct {
	DepDateTime time.Time `json:"DepDateTime"`
	ReturnDateTime time.Time `json:"ReturnDateTime"`
}

type FlightRequest struct {
	PilotID        int       `json:"pilotID"`
	DepDateTime string `json:"depDateTime"`
	ReturnDateTime string `json:"returnDateTime"`
}

func NewPilotInfo(id int, name string, base string, workdays []string) *PilotInfo {
	timeWeekdays, err := stringToWeekday(workdays)
	if err != nil {
		log.Fatal(err)
	}
	pilotInfo := PilotInfo{ID: id, Name: name, Base: base, Workdays: timeWeekdays, Flights: make([]Flight, 0)}
	return &pilotInfo
}

func GetPilotHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idNum, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(pilotInfoMap[idNum])
}

func GetPilotsHandler(w http.ResponseWriter, r *http.Request) {
	pilotNum := len(pilotInfoMap)
	pilotInfos := make([]PilotInfo, 0, pilotNum)
	for i := 1; i <=  pilotNum; i++ {
		pilotInfos = append(pilotInfos, *pilotInfoMap[i])
	}
	json.NewEncoder(w).Encode(pilotInfos)
}

func GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	depDateTime := vars["depDateTime"]
	returnDateTime := vars["returnDateTime"]
	fmt.Println(location + depDateTime + returnDateTime)
}

func PostFlightHandler(w http.ResponseWriter, r *http.Request) {
	var flightRequest FlightRequest
	err := json.NewDecoder(r.Body).Decode(&flightRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(flightRequest)
}