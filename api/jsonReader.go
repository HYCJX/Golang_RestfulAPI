package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"os"
	"time"
)

const dataFile = "crew/crew.json"
var weekdayBitmap = map[time.Weekday]int {
	time.Sunday: 1,
	time.Monday: 2,
	time.Tuesday: 3,
	time.Wednesday: 4,
	time.Thursday : 5,
	time.Friday : 6,
	time.Saturday: 7,
}
var weekdayDic = []time.Weekday {
	time.Sunday,
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
}

func init() {
	// Initialise time-string converter:
	for d := time.Sunday; d <= time.Saturday; d++ {
		name := d.String()
		daysOfWeek[name] = d
		daysOfWeek[name[:3]] = d
	}
	// Read json file:
	crew, err := readJson()
	if err != nil {
		log.Fatal(err)
	}
	// Insert into databse:
	database, _ := sql.Open("sqlite3", "./pilotInfo.db")
	statement, _ := database.Prepare("CREATE TABLE IF Not EXISTS pilots (id INTEGER PRIMARY KEY, name TEXT NOT NULL, base TEXT NOT NULL, workdays SMALLINT ); CREATE TABLE flights (id INTEGER PRIMARY KEY, pilotID INTEGER FOREIGN KEY REFERENCES pilots (id) ON DELETE CASCADE , depDateTime TEXT, returnDateTime TEXT)")
	statement.Exec()
	statement, _ = database.Prepare("INSERT INTO pilots (id, name, base, workdays) VALUES (?,?,?,?)")
	for _, pilot := range crew.Crew {
		timeWeekDays, err := stringToWeekday(pilot.Workdays)
		if err != nil {
			log.Fatal(err)
		}
		statement.Exec(pilot.ID, pilot.Name, pilot.Base, encodeWeekdays(timeWeekDays))
	}
	database.Close()
}

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

func encodeWeekdays(weekdays []time.Weekday) int {
	encodedWeekdays := 0
	for _, weekday := range weekdays {
		encodedWeekdays += int(math.Pow(2, float64(weekdayBitmap[weekday])))
	}
	return encodedWeekdays
}

func decodeWeekdays(encodedWeekdays int) []time.Weekday {
	weekdays := make([]time.Weekday, 0)
	for _, weekday := range weekdayDic {
		if encodedWeekdays & int(math.Pow(2, float64(weekdayBitmap[weekday]))) == 0 {
			weekdays = append(weekdays, weekday)
		}
	}
	return weekdays
}
