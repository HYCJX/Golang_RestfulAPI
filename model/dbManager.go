// Package model directly interacts with database to handle requests from package controller
package model

import (
	"database/sql"
	"github.com/HYCJX/Golang_RestfulAPI/utils"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type PilotInfo struct {
	ID       int      `json:"ID"`
	Name     string   `json:"Name"`
	Base     string   `json:"Base"`
	Workdays []string `json:"Workdays"`
	Flights  []Flight `json:"Flights"`
}

type Flight struct {
	DepDateTime    string `json:"DepDateTime"`
	ReturnDateTime string `json:"ReturnDateTime"`
}

type FlightInfo struct {
	PilotID        int    `json:"PilotID"`
	DepDateTime    string `json:"DepDateTime"`
	ReturnDateTime string `json:"ReturnDateTime"`
}

type PilotIds struct {
	PilotIDs []int `json:"PilotIDs"`
}

type RequestSuccessFlag struct {
	Success bool `json:"Success"`
}

func NewFlight(depDateTime string, returnDateTime string) Flight {
	return Flight{DepDateTime: depDateTime, ReturnDateTime: returnDateTime}
}

func InitDB(dataFile string) {
	// Read json file:
	crew, err := utils.ReadJson(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	// Database initialisation & Insert crew information into database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	// Table 1: pilots.
	statement, _ := database.Prepare("CREATE TABLE IF Not EXISTS pilots (id INTEGER PRIMARY KEY, name TEXT NOT NULL, base TEXT NOT NULL, workdays SMALLINT NOT NULL); ")
	statement.Exec()
	// Table 2: flights.
	statement, _ = database.Prepare(" CREATE TABLE IF NOT EXISTS flights (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, pilotID INTEGER, depDateTime TEXT Not NULL, returnDateTime TEXT NOT NULL, FOREIGN KEY (pilotID) REFERENCES pilots (id) ON DELETE CASCADE);")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO pilots (id, name, base, workdays) VALUES (?,?,?,?)")
	for _, pilot := range crew.Crew {
		timeWeekDays, err := utils.StringToWeekday(pilot.Workdays)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec(pilot.ID, pilot.Name, pilot.Base, utils.EncodeWeekdays(timeWeekDays))
	}
	database.Close()
}

// Handle GET request for a pilot with id.
func GetPilotInfo(id int) (PilotInfo, bool, error) {
	// Open/Close database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	// Initialise results:
	var pilotInfo PilotInfo
	pilotInfo.Flights = make([]Flight, 0)
	hasResult := false
	// Query pilots:
	statement, _ := database.Prepare("SELECT * FROM pilots WHERE id = $1")
	rows, err := statement.Query(id)
	if err != nil {
		return pilotInfo, false, err
	}
	for rows.Next() {
		hasResult = true
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			return pilotInfo, hasResult, err
		}
		pilotInfo.Workdays = utils.DecodeWeekdays(encodedWeekdays)
	}
	// Query flights:
	statement, _ = database.Prepare("SELECT * FROM flights WHERE pilotID = ?")
	rows, err = statement.Query(id)
	if err != nil {
		return pilotInfo, hasResult, err
	}
	for rows.Next() {
		var rowID int
		var pilotID int
		var depDateTime string
		var returnDateTime string
		err := rows.Scan(&rowID, &pilotID, &depDateTime, &returnDateTime)
		if err != nil {
			return pilotInfo, hasResult, err
		}
		pilotInfo.Flights = append(pilotInfo.Flights, NewFlight(depDateTime, returnDateTime))
	}
	return pilotInfo, hasResult, err
}

// Handle GET request for all pilots.
func GetAllPilotsInfo() ([]PilotInfo, error) {
	// Open/Close database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	// Initialise results:
	pilotsInfo := make([]PilotInfo, 0)
	// Query pilots:
	statement, _ := database.Prepare("SELECT * FROM pilots")
	rows, err := statement.Query()
	if err != nil {
		return pilotsInfo, err
	}
	for rows.Next() {
		// Initialise each pilot:
		var pilotInfo PilotInfo
		pilotInfo.Flights = make([]Flight, 0)
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			return pilotsInfo, err
		}
		pilotInfo.Workdays = utils.DecodeWeekdays(encodedWeekdays)
		statement, _ = database.Prepare("SELECT * FROM flights WHERE pilotID = ?")
		// Query flights information for each pilot:
		innerRows, err := statement.Query(pilotInfo.ID)
		if err != nil {
			return pilotsInfo, err
		}
		for innerRows.Next() {
			var id int
			var pilotID int
			var depDateTime string
			var returnDateTime string
			err := innerRows.Scan(&id, &pilotID, &depDateTime, &returnDateTime)
			if err != nil {
				return pilotsInfo, err
			}
			pilotInfo.Flights = append(pilotInfo.Flights, NewFlight(depDateTime, returnDateTime))
		}
		pilotsInfo = append(pilotsInfo, pilotInfo)
	}
	return pilotsInfo, err
}

// Handle GET request for available pilots.
func GetAvailablePilot(location string, depDateTime string, returnDateTime string) (PilotIds, error) {
	// Open/Close database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	// Initialise results:
	var pilotIDs PilotIds
	// Query pilots:
	statement, _ := database.Prepare("SELECT id FROM pilots p_out WHERE base = ? AND workdays & ? <> 0 AND workdays & ? <> 0 AND NOT EXISTS (SELECT p.id, depDateTime, returnDateTime FROM pilots p JOIN flights f on p.id = f.pilotID WHERE p.id = p_out.id AND (depDateTime <= ? AND returnDateTime > ? OR depDateTime > ? and depDateTime < ?))")
	encodedDep, err := utils.EncodeWeekday(depDateTime)
	if err != nil {
		return pilotIDs, err
	}
	encodedRet, err := utils.EncodeWeekday(depDateTime)
	if err != nil {
		return pilotIDs, err
	}
	rows, err := statement.Query(location, encodedDep, encodedRet, depDateTime, depDateTime, depDateTime, returnDateTime)
	if err != nil {
		return pilotIDs, err
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return pilotIDs, err
		}
		pilotIDs.PilotIDs = append(pilotIDs.PilotIDs, id)
	}
	return pilotIDs, err
}

// Handle GET request for flights with pilotID.
func GetFlightInfo(id int) ([]FlightInfo, error) {
	// Open/Close database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	// Initialise results:
	flightInfoList := make([]FlightInfo, 0)
	// Query flights:
	statement, _ := database.Prepare("SELECT * FROM flights WHERE pilotID = $1")
	rows, err := statement.Query(id)
	if err != nil {
		return flightInfoList, err
	}
	for rows.Next() {
		var flightInfo FlightInfo
		var rowID int
		err := rows.Scan(&rowID, &flightInfo.PilotID, &flightInfo.DepDateTime, &flightInfo.ReturnDateTime) // scan contents of the current row into the instance
		if err != nil {
			return flightInfoList, err
		}
		flightInfoList = append(flightInfoList, flightInfo)
	}
	return flightInfoList, err
}

// Handle GET request for all flights.
func GetAllFlightsInfo() ([]FlightInfo, error) {
	// Open/Close database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	// Initialise results:
	flightsInfo := make([]FlightInfo, 0)
	// Query flights:
	statement, _ := database.Prepare("SELECT * FROM flights")
	rows, err := statement.Query()
	if err != nil {
		return flightsInfo, err
	}
	for rows.Next() {
		var flightInfo FlightInfo
		var rowID int
		err := rows.Scan(&rowID, &flightInfo.PilotID, &flightInfo.DepDateTime, &flightInfo.ReturnDateTime) // scan contents of the current row into the instance
		if err != nil {
			return flightsInfo, err
		}
		flightsInfo = append(flightsInfo, flightInfo)
	}
	return flightsInfo, err
}

// Handle POST request for schedule a flight.
func PostFlight(id int, depDateTime string, returnDateTime string) (RequestSuccessFlag, error) {
	// Open/Close database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	// Initialise results:
	var requestSuccessFlag RequestSuccessFlag
	// Conditional insert query:
	statement, _ := database.Prepare("INSERT INTO flights (pilotId, depDateTime, returnDateTime) SELECT ?, ?, ? WHERE EXISTS (SELECT * from pilots WHERE id = ? AND workdays & ? <> 0 AND workdays & ? <> 0) AND NOT EXISTS (SELECT * from flights WHERE pilotID = ? AND (depDateTime <= ? AND returnDateTime > ? OR depDateTime > ? and depDateTime < ?))")
	encodedDep, err := utils.EncodeWeekday(depDateTime)
	if err != nil {
		return requestSuccessFlag, err
	}
	encodedRet, err := utils.EncodeWeekday(depDateTime)
	if err != nil {
		return requestSuccessFlag, err
	}
	rows, err := statement.Exec(id, depDateTime, returnDateTime, id, depDateTime, depDateTime, depDateTime, returnDateTime, id, encodedDep, encodedRet)
	if err != nil {
		return requestSuccessFlag, err
	}
	count, _ := rows.RowsAffected()
	requestSuccessFlag.Success = count != 0
	return requestSuccessFlag, err
}
