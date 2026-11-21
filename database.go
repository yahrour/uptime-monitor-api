package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectToDatabase() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	sqlDb, err := db.DB()

	if err != nil {
		panic("failed to connect database")
	}

	sqlDb.SetMaxOpenConns(1)
	sqlDb.SetMaxIdleConns(1)

	DB = db
}
