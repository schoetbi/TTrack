package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	begin     = kingpin.Command("begin", "Starts a task.")
	beginTask = begin.Arg("task", "Taskname").Required().String()

	end     = kingpin.Command("end", "Ends a new task.")
	endTime = end.Arg("end", "End timestamp '01.02.2019 14:33' or 'now' for current time").String()

	report          = kingpin.Command("report", "Prints a report")
	reportFrom      = report.Arg("from", "From timestamp").String()
	reportTo        = report.Arg("to", "From timestamp").String()
	reportDaily     = report.Flag("daily", "Group times daily").Bool()
	reportByProject = report.Flag("project", "Group by project").Bool()
)

func main() {
	switch kingpin.Parse() {
	case begin.FullCommand():
		BeginTaskHandler(beginTask)
		break
	case end.FullCommand():
		EndOpenTasksHandler(endTime)
	case report.FullCommand():
		ReportHandler(reportFrom, reportTo, *reportDaily, *reportByProject)
	}
}
