package model

import (
	"database/sql"
	"github.com/HYCJX/Golang_RestfulAPI/utils"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type PilotInfo struct {
	ID       int      `json:"ID"`
	Name     string   `json:"Name"`
	Base     string   `json:"Base"`
	Workdays []string `json:"Workdays"`
	Flights  []Flight `json:"Flights"`
}

type Flight struct {
	DepDateTime    time.Time `json:"DepDateTime"`
	ReturnDateTime time.Time `json:"ReturnDateTime"`
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
	statement, _ := database.Prepare("CREATE TABLE IF Not EXISTS pilots (id INTEGER PRIMARY KEY, name TEXT NOT NULL, base TEXT NOT NULL, workdays SMALLINT ); ")
	statement.Exec()
	statement, _ = database.Prepare(" CREATE TABLE IF NOT EXISTS flights (id INTEGER, pilotID INTEGER, depDateTime TEXT, returnDateTime TEXT, PRIMARY KEY (id), FOREIGN KEY (pilotID) REFERENCES pilots (id) ON DELETE CASCADE);")
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
	statement, _ = database.Prepare("SELECT * FROM flights WHERE id = ?")
	rows, err = statement.Query(id)
	for rows.Next() {
		pilotInfo.Flights = append(pilotInfo.Flights)
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
		var encodedWeekdays int
		err := rows.Scan(&pilotInfo.ID, &pilotInfo.Name, &pilotInfo.Base, &encodedWeekdays) // scan contents of the current row into the instance
		if err != nil {
			log.Fatal(err)
		}
		pilotInfo.Workdays = utils.DecodeWeekdays(encodedWeekdays)
		pilotInfos = append(pilotInfos, pilotInfo)
	}
	return pilotInfos
}
