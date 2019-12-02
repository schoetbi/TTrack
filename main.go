package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/olekukonko/tablewriter"
	"github.com/snabb/isoweek"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	begin     = kingpin.Command("begin", "Starts a task.")
	beginTask = begin.Arg("task", "Taskname").Required().String()

	end     = kingpin.Command("end", "Ends a new task.")
	endTime = end.Arg("end", "End timestamp '01.02.2019 14:33' or 'now' for current time").String()

	report      = kingpin.Command("report", "Prints a report")
	reportFrom  = report.Arg("from", "From timestamp").Required().String()
	reportTo    = report.Arg("to", "From timestamp").Required().String()
	reportDaily = report.Flag("daily", "Group times daily").Bool()
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

	endOpenTasks(db, now)

	var task Task
	db.Where(Task{Name: *taskName}).FirstOrCreate(&task)
	fmt.Printf("begin task %s id = %d\n", *taskName, task.ID)
	log := Log{TaskId: task.ID, TimeFrom: now}
	db.Create(&log)
	db.Close()
}

func endOpenTasks(db *gorm.DB, t time.Time) {
	var openLogs []Log
	db.Where("time_to is NULL").Find(&openLogs)
	for _, l := range openLogs {
		var to time.Time = t
		l.TimeTo = &to
		db.Save(&l)
		fmt.Printf("ended log for task %d\n", l.TaskId)
	}
}

func endOpenTasksHandler(endTime *string) {
	var toTime time.Time
	if *endTime == "" || *endTime == "now" {
		toTime = time.Now()
	} else {
		var layout = "2.1.2006"
		parsedTime, err := time.Parse(layout, *endTime)
		if err != nil {
			fmt.Println(err)
			return
		}
		toTime = parsedTime
	}
	var db = getDatabase()
	endOpenTasks(db, toTime)
	fmt.Printf("Finished all open tasks\n")
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

func reportHandler(from *string, to *string, daily bool) {
	layout := "02.01.2006"
	fromTime, err := time.Parse(layout, *from)
	if err != nil {
		fmt.Println(err)
		return
	}
	toTime, err := time.Parse(layout, *to)
	if err != nil {
		fmt.Println(err)
		return
	}
	var db = getDatabase()
	defer db.Close()
	fmt.Printf("Report from:%s to:%s\n", *from, *to)
	type Result struct {
		TaskId       uint
		Day          int
		Name         string
		TotalSeconds float64
	}
	table := tablewriter.NewWriter(os.Stdout)
	var results []Result
	if daily {
		table.SetHeader([]string{"Day", "Task", "Time [h]"})
		db.Table("logs").
			Select("cast(round(julianday(logs.time_from)) as int) as day, logs.task_id, tasks.name, sum((julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0) as total_seconds").
			Joins("join tasks on tasks.id = logs.task_id").
			Where("time_from > ? and time_to < ?", fromTime, toTime).
			Group("day, logs.task_id").
			Find(&results)
		last_day := 0
		for _, r := range results {
			if last_day != r.Day {
				y, month, day := isoweek.JulianToDate(r.Day)
				date := time.Date(y, month, day, 0, 0, 0, 0, time.Local)

				row := []string{date.Format("02.01.2006"), r.Name, fmt.Sprintf("%f", r.TotalSeconds/60/60)}
				table.Append(row)
				last_day = r.Day
			} else {
				row := []string{" ", r.Name, fmt.Sprintf("%f", r.TotalSeconds/60/60)}
				table.Append(row)
			}
		}
		table.Render()
	} else {
		table.SetHeader([]string{"Task", "Time [h]"})
		db.Table("logs").
			Select("logs.task_id, tasks.name, sum((julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0) as total_seconds").
			Joins("join tasks on tasks.id = logs.task_id").
			Where("time_from > ? and time_to < ?", fromTime, toTime).
			Group("logs.task_id").
			Find(&results)
		for _, r := range results {
			row := []string{r.Name, fmt.Sprintf("%f", r.TotalSeconds/60/60)}
			table.Append(row)
		}
		table.Render()
	}
}

func main() {
	switch kingpin.Parse() {
	case begin.FullCommand():
		beginTaskHandler(beginTask)
		break
	case end.FullCommand():
		endOpenTasksHandler(endTime)
	case report.FullCommand():
		reportHandler(reportFrom, reportTo, *reportDaily)
	}
}
