package api

import (
	"encoding/json"
	"os"
	"time"
)

const dataFile = "crew/crew.json"

type Crew struct {
	ID int `json:"ID"`
	Name string `json:"Name"`
	Base string `json:"Base"`
	Workdays []time.Weekday `json:"Workdays"`
}

func ReadJson() (*Crew, error) {
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