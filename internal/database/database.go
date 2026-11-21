package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{
		TranslateError: true,
	})

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
