package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lutfiharidha/pokemon/app/types"
	"github.com/lutfiharidha/pokemon/db"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var dataLog []string

type LogModel interface {
	Init()
	SaveLog(req types.Logging)
	IntervalLog(req string)
	SpecificLog(req string)
	SpecificIntervalLog(from string, to string)
}

type logModel struct {
	db *gorm.DB
}

func NewLogModels() LogModel {
	return &logModel{}
}

func (m *logModel) Init() {
	db := db.NewSQL().SetupDatabaseConnection()
	m.db = db //calling database connection
}

// function to save log
func (m *logModel) SaveLog(data types.Logging) {
	jData, _ := json.Marshal(data.DataLog)
	req := types.Log{
		DataLog: datatypes.JSON([]byte(jData)),
		Winner:  data.Winner,
	}
	tx := m.db.Create(&req)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
}

// function to display logs based on time interval
func (m *logModel) IntervalLog(req string) {
	dmy := string(req[len(req)-1])
	numT, _ := strconv.Atoi(string(req[:len(req)-1]))
	now := time.Now()
	after := time.Now()
	switch strings.ToLower(dmy) {
	case "m":
		after = now.AddDate(0, -numT, 0)
	case "y":
		after = now.AddDate(-numT, 0, 0)
	case "d":
		after = now.AddDate(0, 0, -numT)
	default:
		log.Fatal("Example input 1d [1 day]")
	}
	var logData []types.Log
	tx1 := m.db.Where("created_at BETWEEN ? and ?", after, now).Find(&logData)
	if tx1.Error != nil {
		log.Fatal(tx1.Error)
	}
	if len(logData) <= 0 {
		log.Println("NO DATA in", req)
		os.Exit(0)
	}
	date := []string{
		after.Format("2006-01-02"),
		now.Format("2006-01-02"),
	}
	saveLogToFile(logData, date) //save log into file
}

// function to display logs specific date
func (m *logModel) SpecificLog(req string) {
	var logData []types.Log

	tx1 := m.db.Where("date(created_at) = ?", req).Find(&logData)
	if tx1.Error != nil {
		log.Fatal(tx1.Error)
	}
	if len(logData) <= 0 {
		log.Println("NO DATA in", req)
		os.Exit(0)
	}
	date := []string{
		req,
	}
	saveLogToFile(logData, date) //save log into file

}

// function to display logs based on specific time interval
func (m *logModel) SpecificIntervalLog(from string, to string) {
	var logData []types.Log

	tx1 := m.db.Where("date(created_at) BETWEEN ? AND ?", from, to).Find(&logData)
	if tx1.Error != nil {
		log.Fatal(tx1.Error)
	}
	if len(logData) <= 0 {
		log.Println("NO DATA from ", from, "-", to)
		os.Exit(0)
	}
	date := []string{
		from,
		to,
	}
	saveLogToFile(logData, date) //save log into file

}

// function to print log
func showLog(str string) {
	fmt.Printf(str)
	goToDB := strings.ReplaceAll(str, "\n", "")
	dataLog = append(dataLog, goToDB)
}

// function to save log into file
func saveLogToFile(datas []types.Log, time []string) {

	f, err := os.Create("./battle.log")
	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	if len(time) > 1 {
		f.WriteString("Data Log from " + time[0] + " to " + time[1] + "\n\n")
	} else {
		f.WriteString("Data Log " + time[0] + "\n\n")

	}
	for i, data := range datas {
		if i != 0 || i != len(datas)-1 {
			f.WriteString("ANOTHER GAME \n\n")
		}
		str := strings.Split(string(data.DataLog), ",")

		for _, v := range str {
			first := strings.ReplaceAll(v, `[`, "")
			second := strings.ReplaceAll(first, `]`, "")
			third := strings.ReplaceAll(second, `"`, "")
			if strings.Contains(third, "THE WINNER") {
				f.WriteString(third + "\n\n")
			} else {
				f.WriteString(third + "\n")
			}
		}
	}
	fmt.Println("DONE save to ./battle.log")
}
