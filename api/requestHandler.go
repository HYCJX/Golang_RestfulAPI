package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PilotInfo struct {
	ID       int            `json:"ID"`
	Name     string         `json:"Name"`
	Base     string         `json:"Base"`
	Workdays []time.Weekday `json:"Workdays"`
	Flights  []Flight       `json:"Flights"`
}

type Flight struct {
	DepDateTime    time.Time `json:"DepDateTime"`
	ReturnDateTime time.Time `json:"ReturnDateTime"`
}

type FlightRequest struct {
	PilotID        int    `json:"pilotID"`
	DepDateTime    string `json:"depDateTime"`
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
	//id := mux.Vars(r)["id"]
	//idNum, err := strconv.Atoi(id)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//json.NewEncoder(w).Encode(pilotInfoMap[idNum])
}

func GetPilotsHandler(w http.ResponseWriter, r *http.Request) {
	//pilotNum := len(pilotInfoMap)
	//pilotInfos := make([]PilotInfo, 0, pilotNum)
	//for i := 1; i <= pilotNum; i++ {
	//	pilotInfos = append(pilotInfos, *pilotInfoMap[i])
	//}
	//json.NewEncoder(w).Encode(pilotInfos)
}

func GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	depDateTime, err1 := stringToTime(vars["depDateTime"])
	if err1 != nil {
		log.Fatal(err1)
	}
	returnDateTime, err2 := stringToTime(vars["returnDateTime"])
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(location + depDateTime.String() + returnDateTime.String())
}

func PostFlightHandler(w http.ResponseWriter, r *http.Request) {
	var flightRequest FlightRequest
	err := json.NewDecoder(r.Body).Decode(&flightRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pilotId := flightRequest.PilotID
	depDateTime, err1 := stringToTime(flightRequest.DepDateTime)
	if err1 != nil {
		log.Fatal(err1)
	}
	returnDateTime, err2 := stringToTime(flightRequest.ReturnDateTime)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(strconv.Itoa(pilotId) + depDateTime.String() + returnDateTime.String())
}
