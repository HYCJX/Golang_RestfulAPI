// Package controller contains functions that handles http requests by connecting to package model
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
	pilotInfo, hasResult, err := model.GetPilotInfo(idNum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	if !hasResult {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		json.NewEncoder(w).Encode(pilotInfo)
	}
}

func GetPilotsHandler(w http.ResponseWriter, r *http.Request) {
	pilotsInfo, err := model.GetAllPilotsInfo()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(pilotsInfo)
}

func GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	depDateTime, err := utils.FormatTimeString(vars["depDateTime"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	returnDateTime, err := utils.FormatTimeString(vars["returnDateTime"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	pilotIds, err := model.GetAvailablePilot(location, depDateTime, returnDateTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(pilotIds)
}

func GetFlightHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idNum, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	flightInfo, err := model.GetFlightInfo(idNum)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(flightInfo)
}

func GetFlightsHandler(w http.ResponseWriter, r *http.Request) {
	flightsInfo, err := model.GetAllFlightsInfo()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
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
	depDateTime, err := utils.FormatTimeString(flightRequest.DepDateTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	returnDateTime, err := utils.FormatTimeString(flightRequest.ReturnDateTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	flag, err := model.PostFlight(pilotId, depDateTime, returnDateTime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(flag)
}
