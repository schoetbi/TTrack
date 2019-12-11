package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/mitchellh/go-homedir"
)

type Task struct {
	ID      uint
	Project string
	Name    string
	Logs    []Log
}

type Log struct {
	ID       uint
	TimeFrom time.Time
	TimeTo   *time.Time
	TaskId   uint
}

func EndOpenTasks(db *gorm.DB, t time.Time) {
	var openLogs []Log
	db.Where("time_to is NULL").Find(&openLogs)
	for _, l := range openLogs {
		var to time.Time = t
		l.TimeTo = &to
		db.Save(&l)
		duration := to.Sub(l.TimeFrom)
		fmt.Printf("ended log for task %d %s\n", l.TaskId, FormatTime(duration.Seconds()))
	}
}

func GetDatabase() *gorm.DB {
	homeDir, errHomeDir := homedir.Dir()
	if errHomeDir != nil {
		panic("Unable to get home directory")
	}

	ttrackPath := path.Join(homeDir, "ttrack")
	os.MkdirAll(ttrackPath, os.ModePerm)
	fullPath := path.Join(ttrackPath, "ttrack.db")
	fmt.Printf("Using database at %s\n", fullPath)
	db, err := gorm.Open("sqlite3", fullPath)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Task{})
	return db
}
