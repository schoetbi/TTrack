package main

import (
	"fmt"

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

func main() {
	switch kingpin.Parse() {
	case begin.FullCommand():
		fmt.Printf("begin task %s\n", *beginTask)
		break
	case end.FullCommand():
		fmt.Printf("end\n")
	case report.FullCommand():
		fmt.Printf("Report %s-%s\n", *reportFrom, *reportTo)
	}
}
