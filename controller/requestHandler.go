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
	depDateTime, err1 := utils.StringToTime(vars["depDateTime"])
	if err1 != nil {
		log.Fatal(err1)
	}
	returnDateTime, err2 := utils.StringToTime(vars["returnDateTime"])
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
	depDateTime, err1 := utils.StringToTime(flightRequest.DepDateTime)
	if err1 != nil {
		log.Fatal(err1)
	}
	returnDateTime, err2 := utils.StringToTime(flightRequest.ReturnDateTime)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(strconv.Itoa(pilotId) + depDateTime.String() + returnDateTime.String())
}
