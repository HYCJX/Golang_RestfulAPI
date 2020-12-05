package api

import (
	"encoding/json"
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
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var crew *Crew
	err = json.NewDecoder(file).Decode(&crew)
	return crew, err
}
