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
	whereClause := getWhereClause(from, to)
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
				Where(whereClause).
				Group("day, project").
				Find(&results)
			last_day := 0
			for _, r := range results {
				if last_day != r.Day {
					y, month, day := isoweek.JulianToDate(r.Day)
					date := time.Date(y, month, day, 0, 0, 0, 0, time.Local)
					row := []string{date.Format("02.01.2006"), r.Project, FormatTime(r.TotalSeconds)}
					table.Append(row)
					last_day = r.Day
				} else {
					row := []string{" ", r.Project, FormatTime(r.TotalSeconds)}
					table.Append(row)
				}
			}
			table.Render()
		} else {
			table.SetHeader([]string{"Project", "Time [h]"})
			db.Table("logs").
				Select("tasks.project as project, sum((julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0) as total_seconds").
				Joins("join tasks on tasks.id = logs.task_id").
				Where(whereClause).
				Group("project").
				Find(&results)
			for _, r := range results {
				row := []string{r.Project, FormatTime(r.TotalSeconds)}
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
				Where(whereClause).
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
					row := []string{" ", r.Name, FormatTime(r.TotalSeconds)}
					table.Append(row)
				}
			}
			table.Render()
		} else {
			table.SetHeader([]string{"Task", "Time [h]"})
			db.Table("logs").
				Select("logs.task_id, tasks.name, sum((julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0) as total_seconds").
				Joins("join tasks on tasks.id = logs.task_id").
				Where(whereClause).
				Group("logs.task_id").
				Find(&results)
			for _, r := range results {
				row := []string{r.Name, FormatTime(r.TotalSeconds)}
				table.Append(row)
			}
			table.Render()
		}
	}
}

func getWhereClause(from *string, to *string) string {
	sqliteDateFormat := "2006-01-02 15:04:05"
	var fromTime time.Time
	var toTime time.Time
	//today
	if strings.HasPrefix(*from, "t") {
		// today
		now := time.Now()
		fromTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	} else if strings.HasPrefix(*from, "y") {
		// yesterday
		yesterday := time.Now().AddDate(0, 0, -1)
		fromTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	} else if strings.HasPrefix(*from, "m") {
		// current month
		now := time.Now()
		currentYear, currentMonth, _ := now.Date()
		fromTime = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
	} else if strings.HasPrefix(*from, "pm") {
		// previous month
		now := time.Now()
		lastMonth := now.AddDate(0, -1, 0)
		lastMonthYear, lastMonthMonth, _ := lastMonth.Date()
		fromTime = time.Date(lastMonthYear, lastMonthMonth, 1, 0, 0, 0, 0, lastMonth.Location())

		currentYear, currentMonth, _ := now.Date()
		toTime = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
		return fmt.Sprintf("time_from >= '%s' and time_to <= '%s'", fromTime.Format(sqliteDateFormat), toTime.Format(sqliteDateFormat))
	} else if *from != "" {
		t, err := ParseDateTime(*from)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		fromTime = t
	}
	if *to != "" {
		t, err := ParseDateTime(*to)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		toTime = t
	}

	var whereClause string
	if *from == "" && *to != "" {
		whereClause = "true"
	} else if *from != "" && *to == "" {
		whereClause = fmt.Sprintf("time_from >= '%s'", fromTime.Format(sqliteDateFormat))
	} else if *from == "" && *to != "" {
		whereClause = fmt.Sprintf("time_to <= '%s'", toTime.Format(sqliteDateFormat))
	} else if *from != "" && *to != "" {
		whereClause = fmt.Sprintf("time_from >= '%s' and time_to <= '%s'", fromTime.Format(sqliteDateFormat), toTime.Format(sqliteDateFormat))
	}
	return whereClause
}

func ListHandler(from *string, to *string) {
	whereClause := getWhereClause(from, to)
	fmt.Println(whereClause)
	var db = GetDatabase()
	defer db.Close()
	fmt.Printf("List from:%s to:%s\n", *from, *to)
	type Result struct {
		TaskId       uint
		Name         string
		TimeTo       time.Time
		TimeFrom     time.Time
		Project      string
		TotalSeconds float64
	}
	table := tablewriter.NewWriter(os.Stdout)
	var results []Result
	table.SetHeader([]string{"Task", "From", "To", "Time [h]"})
	db.Table("logs").
		Select("logs.task_id, tasks.name, logs.time_from, logs.time_to, (julianday(logs.time_to) - julianday(logs.time_from)) * 86400.0 as total_seconds").
		Joins("join tasks on tasks.id = logs.task_id").
		Where(whereClause).
		Order("logs.time_from").
		Find(&results)
	dateTimeFormat := "02.01.2006 15:04:05"
	zeroTime := time.Time{}
	for _, r := range results {
		var row []string
		if r.TimeTo == zeroTime {
			tempDiff := time.Now().Sub(r.TimeFrom)
			row = []string{r.Name, r.TimeFrom.Format(dateTimeFormat), "", FormatTime(tempDiff.Seconds())}
		} else {
			row = []string{r.Name, r.TimeFrom.Format(dateTimeFormat), r.TimeTo.Format(dateTimeFormat), FormatTime(r.TotalSeconds)}
		}

		table.Append(row)
	}
	table.Render()
}

func FormatTime(timeInSeconds float64) string {
	return fmt.Sprintf("%f (%.1f min)", timeInSeconds/60/60, timeInSeconds/60.0)
}

func BeginTaskHandler(taskName *string, startTime *string) {
	var db = GetDatabase()

	var now time.Time
	if startTime != nil && *startTime != "" {
		start, err := ParseDateTime(*startTime)
		if err != nil {
			return
		}
		if start.Year() == 0 {
			today := time.Now()
			now = time.Date(today.Year(), today.Month(), today.Day(), start.Hour(), start.Minute(), start.Second(), 0, time.Local)
		} else {
			now = start
		}

	} else {
		now = time.Now()
	}

	EndOpenTasks(db, now)
	var splitted = strings.Split(*taskName, "/")
	var task Task
	if len(splitted) == 2 {
		db.Where(Task{Project: splitted[0], Name: *taskName}).FirstOrCreate(&task)
	} else {
		db.Where(Task{Name: *taskName}).FirstOrCreate(&task)
	}

	fmt.Printf("begin task %s\n", *taskName)
	log := Log{TaskId: task.ID, TimeFrom: now}
	db.Create(&log)
	db.Close()
}

func EndOpenTasksHandler(endTime *string) {
	var toTime time.Time
	if *endTime == "" || *endTime == "now" {
		toTime = time.Now()
	} else {
		var layout = "2.1.2006 15:04"
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
