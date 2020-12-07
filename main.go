package main

import (
	"github.com/HYCJX/Golang_RestfulAPI/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/pilots/availability", controller.GetAvailabilityHandler).
		Queries("location", "{location}", "depDateTime", "{depDateTime}", "returnDateTime", "{returnDateTime}").
		Methods("GET")
	router.HandleFunc("/pilots/{id:[0-9]+}", controller.GetPilotHandler).Methods("GET")
	router.HandleFunc("/pilots", controller.GetPilotsHandler).Methods("GET")
	router.HandleFunc("/flights/{id:[0-9]+}", controller.GetFlightHandler).Methods("GET")
	router.HandleFunc("/flights", controller.GetFlightsHandler).Methods("GET")
	router.HandleFunc("/flights", controller.PostFlightHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
