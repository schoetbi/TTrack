package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/snabb/isoweek"
)

func ReportHandler(from *string, to *string, daily bool, byProject bool) {
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
	var db = GetDatabase()
	defer db.Close()
	fmt.Printf("Report from:%s to:%s\n", *from, *to)
	type Result struct {
		TaskId       uint
		Day          int
		Name         string
		Project      string
		TotalSeconds float64
	}
	table := tablewriter.NewWriter(os.Stdout)
	var results []Result
	if byProject {
		if daily {
			table.SetHeader([]string{"Day", "Project", "Time [h]"})
			db.Table("logs").
				Select("cast(round(julianday(logs.time_from)) as int) as day, tasks.project as project, sum((julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0) as total_seconds").
				Joins("join tasks on tasks.id = logs.task_id").
				Where("time_from > ? and time_to < ?", fromTime, toTime).
				Group("day, project").
				Find(&results)
			last_day := 0
			for _, r := range results {
				if last_day != r.Day {
					y, month, day := isoweek.JulianToDate(r.Day)
					date := time.Date(y, month, day, 0, 0, 0, 0, time.Local)

					row := []string{date.Format("02.01.2006"), r.Project, fmt.Sprintf("%f", r.TotalSeconds/60/60)}
					table.Append(row)
					last_day = r.Day
				} else {
					row := []string{" ", r.Project, fmt.Sprintf("%f", r.TotalSeconds/60/60)}
					table.Append(row)
				}
			}
			table.Render()
		} else {
			table.SetHeader([]string{"Project", "Time [h]"})
			db.Table("logs").
				Select("tasks.project as project, sum((julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0) as total_seconds").
				Joins("join tasks on tasks.id = logs.task_id").
				Where("time_from > ? and time_to < ? and tasks.project is not NULL", fromTime, toTime).
				Group("project").
				Find(&results)
			for _, r := range results {
				row := []string{r.Project, fmt.Sprintf("%f", r.TotalSeconds/60/60)}
				table.Append(row)
			}
			table.Render()
		}
	} else {
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
}

func BeginTaskHandler(taskName *string) {
	var now = time.Now()
	var db = GetDatabase()

	EndOpenTasks(db, now)
	var splitted = strings.Split(*taskName, "/")
	var task Task
	if len(splitted) == 2 {
		db.Where(Task{Project: splitted[0], Name: *taskName}).FirstOrCreate(&task)
	} else {
		db.Where(Task{Name: *taskName}).FirstOrCreate(&task)
	}

	fmt.Printf("begin task %s id = %d\n", *taskName, task.ID)
	log := Log{TaskId: task.ID, TimeFrom: now}
	db.Create(&log)
	db.Close()
}

func EndOpenTasksHandler(endTime *string) {
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
	var db = GetDatabase()
	EndOpenTasks(db, toTime)
	fmt.Printf("Finished all open tasks\n")
}
