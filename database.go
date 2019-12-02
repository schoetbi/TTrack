package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
		fmt.Printf("ended log for task %d\n", l.TaskId)
	}
}

func GetDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "ttrack.db")
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Task{})
	return db
}
