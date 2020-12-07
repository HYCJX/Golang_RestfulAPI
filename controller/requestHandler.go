package controller

import (
	"encoding/json"
	"github.com/HYCJX/Golang_RestfulAPI/model"
	"github.com/HYCJX/Golang_RestfulAPI/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

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
	pilotInfo := model.GetPilotInfo(idNum)
	json.NewEncoder(w).Encode(pilotInfo)
}

func GetPilotsHandler(w http.ResponseWriter, r *http.Request) {
	pilotsInfo := model.GetAllPilotsInfo()
	json.NewEncoder(w).Encode(pilotsInfo)
}

func GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	depDateTime := utils.FormatTimeString(vars["depDateTime"])
	returnDateTime := utils.FormatTimeString(vars["returnDateTime"])
	pilotIds := model.GetAvailablePilot(location, depDateTime, returnDateTime)
	json.NewEncoder(w).Encode(pilotIds)
}

func GetFlightHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idNum, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	flightInfo := model.GetFlightInfo(idNum)
	json.NewEncoder(w).Encode(flightInfo)
}

func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	flightsInfo := model.GetAllFlightsInfo()
	json.NewEncoder(w).Encode(flightsInfo)
}

func PostFlightHandler(w http.ResponseWriter, r *http.Request) {
	var flightRequest FlightRequest
	err := json.NewDecoder(r.Body).Decode(&flightRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pilotId := flightRequest.PilotID
	depDateTime := utils.FormatTimeString(flightRequest.DepDateTime)
	returnDateTime := utils.FormatTimeString(flightRequest.ReturnDateTime)
	flag := model.PostFlight(pilotId, depDateTime, returnDateTime)
	json.NewEncoder(w).Encode(flag)
}
