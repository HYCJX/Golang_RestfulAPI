package api

import (
	"database/sql"
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
	ID       int      `json:"ID"`
	Name     string   `json:"Name"`
	Base     string   `json:"Base"`
	Workdays []string `json:"Workdays"`
	Flights  []Flight `json:"Flights"`
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

func GetPilotHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idNum, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("SELECT * FROM pilots WHERE id = $1")
	rows, err := statement.Query(idNum)
	if err != nil {
		log.Fatal(err)
	}
	var pilotInfo PilotInfo
	for rows.Next() {
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		pilotInfo.Workdays = decodeWeekdays(encodedWeekdays)
	}
	json.NewEncoder(w).Encode(pilotInfo)
}

func GetPilotsHandler(w http.ResponseWriter, r *http.Request) {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("SELECT * FROM pilots")
	rows, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	pilotInfos := make([]PilotInfo, 0)
	for rows.Next() {
		var pilotInfo PilotInfo
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		pilotInfo.Workdays = decodeWeekdays(encodedWeekdays)
		pilotInfos = append(pilotInfos, pilotInfo)
	}
	json.NewEncoder(w).Encode(pilotInfos)
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
