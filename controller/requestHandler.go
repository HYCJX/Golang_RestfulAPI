package controller

import (
	"encoding/json"
	"fmt"
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
	pilotInfos := model.GetAllPilotsInfo()
	json.NewEncoder(w).Encode(pilotInfos)
}

func GetAvailabilityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	location := vars["location"]
	depDateTime := utils.FormatTimeString(vars["depDateTime"])

	returnDateTime := utils.FormatTimeString(vars["returnDateTime"])

	fmt.Println(location + depDateTime + returnDateTime)
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
	model.PostFlight(pilotId, depDateTime, returnDateTime)
}
