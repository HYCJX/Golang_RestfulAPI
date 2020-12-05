package main

import (
	"github.com/HYCJX/Golang_RestfulAPI/api"
)

func main() {
	// Initialise map:
	pilotInfoMap := make(map[int]*api.PilotInfo)
	api.GetPilotInfo(pilotInfoMap)
}
