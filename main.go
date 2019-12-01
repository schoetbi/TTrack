package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	begin     = kingpin.Command("begin", "Starts a task.")
	beginTask = begin.Arg("task", "Taskname").Required().String()

	end = kingpin.Command("end", "Ends a new task.")

	report     = kingpin.Command("report", "Prints a report")
	reportFrom = report.Flag("from", "From timestamp").Required().String()
	reportTo   = report.Flag("to", "From timestamp").Required().String()
)

type Task struct {
	ID   uint
	Name string
	Logs []Log
}

type Log struct {
	ID       uint
	TimeFrom time.Time
	TimeTo   *time.Time
	TaskId   uint
}

func beginTaskHandler(taskName *string) {
	var now = time.Now()
	var db = getDatabase()

	// are there open logs?
	var openLogs []Log
	db.Where("time_to is NULL").Find(&openLogs)
	for _, l := range openLogs {
		var to time.Time = now
		l.TimeTo = &to
		db.Save(&l)
		fmt.Printf("ended log for task %d\n", l.TaskId)
	}

	var task Task
	db.Where(Task{Name: *taskName}).FirstOrCreate(&task)
	fmt.Printf("begin task %s id = %d\n", *taskName, task.ID)
	log := Log{TaskId: task.ID, TimeFrom: now}
	db.Create(&log)
	db.Close()
}

func getDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "ttrack.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Task{})
	return db
}

func main() {

	switch kingpin.Parse() {
	case begin.FullCommand():
		beginTaskHandler(beginTask)
		break
	case end.FullCommand():
		fmt.Printf("end\n")
	case report.FullCommand():
		layout := "2.1.2006"
		fromTime, err := time.Parse(layout, *reportFrom)
		if err != nil {
			fmt.Println(err)
		}
		toTime, _ := time.Parse(layout, *reportTo)
		fmt.Printf("Report %s-%s\n", fromTime, toTime)
	}
}
