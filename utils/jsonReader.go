package utils

import (
	"encoding/json"
	"os"
)

type Crew struct {
	Crew []CrewItem `json:"Crew"`
}

type CrewItem struct {
	ID       int      `json:"ID"`
	Name     string   `json:"Name"`
	Base     string   `json:"Base"`
	Workdays []string `json:"Workdays"`
}

func ReadJson(dataFile string) (*Crew, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var crew *Crew
	err = json.NewDecoder(file).Decode(&crew)
	return crew, err
}
