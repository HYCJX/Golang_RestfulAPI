package main

import (
	"github.com/HYCJX/Golang_RestfulAPI/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/pilots/{id:[0-9]+}", api.GetPilotHandler).Methods("GET")
	router.HandleFunc("/pilots/availability", api.GetAvailabilityHandler).
		Queries("location", "{location}", "depDateTime", "{depDateTime}", "returnDateTime", "{returnDateTime}").
		Methods("GET")
	router.HandleFunc("/pilots", api.GetPilotsHandler).Methods("GET")
	router.HandleFunc("/flights", api.PostFlightHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
