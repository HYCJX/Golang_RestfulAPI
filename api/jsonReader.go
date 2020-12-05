package api

import (
	"encoding/json"
	"log"
	"os"
)

const dataFile = "crew/crew.json"

type Crew struct {
	Crew []CrewItem `json:"Crew"`
}

type CrewItem struct {
	ID       int      `json:"ID"`
	Name     string   `json:"Name"`
	Base     string   `json:"Base"`
	Workdays []string `json:"Workdays"`
}

func readJson() (*Crew, error) {
	// Open the file.
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}

	// Schedule the file to be closed once
	// the function returns.
	defer file.Close()

	// Decode the file into a slice of pointers
	// to Feed values.
	var crew *Crew
	err = json.NewDecoder(file).Decode(&crew)

	// We don't need to check for errors, the caller can do this.
	return crew, err
}

func GetPilotInfo(pilotInfoMap map[int]*PilotInfo) {
	// Read json file:
	crew, err := readJson()
	if err != nil {
		log.Fatal(err)
	}
	for _, pilot := range crew.Crew {
		pilotInfoMap[pilot.ID] = newPilotInfo(pilot.Name, pilot.Base, pilot.Workdays)
	}
}
