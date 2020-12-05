package main

import (
	"github.com/HYCJX/Golang_RestfulAPI/api"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/pilots/{id}", api.GetPilotHandler).Methods("GET")
	router.HandleFunc("/pilots", api.GetPilotsHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}