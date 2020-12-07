package model

import (
	"database/sql"
	"fmt"
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

func init() {
	// Read json file:
	const dataFile = "crew/crew.json"
	crew, err := utils.ReadJson(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	// Database initialisation & Insert crew information into database:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	statement, _ := database.Prepare("CREATE TABLE IF Not EXISTS pilots (id INTEGER PRIMARY KEY, name TEXT NOT NULL, base TEXT NOT NULL, workdays SMALLINT NOT NULL); ")
	statement.Exec()
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

func GetPilotInfo(id int) PilotInfo {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("SELECT * FROM pilots WHERE id = $1")
	rows, err := statement.Query(id)
	if err != nil {
		fmt.Println("Here")
		log.Fatal(err)
	}

	var pilotInfo PilotInfo
	pilotInfo.Flights = make([]Flight, 0)
	for rows.Next() {
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		pilotInfo.Workdays = utils.DecodeWeekdays(encodedWeekdays)
	}
	statement, _ = database.Prepare("SELECT * FROM flights WHERE pilotID = ?")
	rows, err = statement.Query(id)
	if (err != nil) {
		log.Fatal(err)
	}
	for rows.Next() {
		var rowID int
		var pilotID int
		var depDateTime string
		var returnDateTime string
		err := rows.Scan(&rowID, &pilotID, &depDateTime, &returnDateTime)
		if err != nil {
			log.Fatal(err)
		}
		pilotInfo.Flights = append(pilotInfo.Flights, NewFlight(depDateTime, returnDateTime))
	}
	return pilotInfo
}

func GetAllPilotsInfo() []PilotInfo {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("SELECT * FROM pilots")
	rows, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	pilotsInfo := make([]PilotInfo, 0)
	for rows.Next() {
		var pilotInfo PilotInfo
		pilotInfo.Flights = make([]Flight, 0)
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		pilotInfo.Workdays = utils.DecodeWeekdays(encodedWeekdays)
		statement, _ = database.Prepare("SELECT * FROM flights WHERE pilotID = ?")
		innerRows, err := statement.Query(pilotInfo.ID)
		for innerRows.Next() {
			var id int
			var pilotID int
			var depDateTime string
			var returnDateTime string
			err := innerRows.Scan(&id, &pilotID, &depDateTime, &returnDateTime)
			if err != nil {
				log.Fatal(err)
			}
			pilotInfo.Flights = append(pilotInfo.Flights, NewFlight(depDateTime, returnDateTime))
		}
		pilotsInfo = append(pilotsInfo, pilotInfo)
	}
	return pilotsInfo
}

func GetFlightInfo(id int) []FlightInfo {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("SELECT * FROM flights WHERE pilotID = $1")
	rows, err := statement.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	flightInfoList := make([]FlightInfo, 0)
	for rows.Next() {
		var flightInfo FlightInfo
		var rowID int
		err := rows.Scan(&rowID, &flightInfo.PilotID, &flightInfo.DepDateTime, &flightInfo.ReturnDateTime) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		flightInfoList = append(flightInfoList, flightInfo)
	}
	return flightInfoList
}

func GetAllFlightsInfo() []FlightInfo {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("SELECT * FROM pilots")
	rows, err := statement.Query()
	if err != nil {
		log.Fatal(err)
	}
	flightsInfo := make([]FlightInfo, 0)
	for rows.Next() {
		var flightInfo FlightInfo
		var rowID int
		err := rows.Scan(&rowID, &flightInfo.PilotID, &flightInfo.DepDateTime, &flightInfo.ReturnDateTime) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		flightsInfo = append(flightsInfo, flightInfo)
	}
	return flightsInfo
}

func GetAvailablePilot(location string, depDateTime string, returnDateTime string) PilotIds {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	var pilotIDs PilotIds
	statement, _ := database.Prepare("SELECT id FROM pilots p_out WHERE base = ? AND NOT EXISTS (SELECT p.id, depDateTime, returnDateTime FROM pilots p JOIN flights f on p.id = f.pilotID WHERE p.id = p_out.id AND (depDateTime <= ? AND returnDateTime > ? OR depDateTime > ? and depDateTime < ?))")
	rows, err := statement.Query(location, depDateTime, depDateTime, depDateTime, returnDateTime)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		pilotIDs.PilotIDs = append(pilotIDs.PilotIDs, id)
	}
	return pilotIDs
}

func PostFlight(id int, depDateTime string, returnDateTime string) RequestSuccessFlag {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("INSERT INTO flights (pilotId, depDateTime, returnDateTime) SELECT ?, ?, ? WHERE NOT EXISTS (SELECT * from flights WHERE pilotID = ? AND (depDateTime <= ? AND returnDateTime > ? OR depDateTime > ? and depDateTime < ?))")
	rows, err := statement.Exec(id, depDateTime, returnDateTime, id, depDateTime, depDateTime, depDateTime, returnDateTime)
	if err != nil {
		log.Fatal(err)
	}
	count, _ := rows.RowsAffected()
	var requestSuccessFlag RequestSuccessFlag
	requestSuccessFlag.Success = count != 0
	return requestSuccessFlag
}
