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
	for rows.Next() {
		var rowID int
		var pilotId int
		var depDateTime string
		var returnDateTime string
		err := rows.Scan(&rowID, &pilotId, &depDateTime, &returnDateTime)
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
	pilotInfos := make([]PilotInfo, 0)
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
		rows, err = statement.Query(pilotInfo.ID)
		for rows.Next() {
			var id int
			var pilotId int
			var depDateTime string
			var returnDateTime string
			err := rows.Scan(&id, &pilotId, &depDateTime, &returnDateTime)
			if err != nil {
				log.Fatal(err)
			}
			pilotInfo.Flights = append(pilotInfo.Flights, NewFlight(depDateTime, returnDateTime))
		}
		pilotInfos = append(pilotInfos, pilotInfo)
	}
	return pilotInfos
}

func PostFlight(id int, depDateTime string, returnDateTime string) {
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	defer database.Close()
	statement, _ := database.Prepare("INSERT INTO flights (pilotId, depDateTime, returnDateTime) VALUES (?,?,?)")
	statement.Exec(id, depDateTime, returnDateTime)
	statement, _ = database.Prepare("SELECT * FROM flights WHERE depDateTime > ?")
	rows, _ := statement.Query("2025-08-01T10:00:00Z")
	for rows.Next() {
		var id int
		var pilotId int
		var depDateTime string
		var returnDateTime string
		err := rows.Scan(&id, &pilotId, &depDateTime, &returnDateTime)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(depDateTime)
	}
}
