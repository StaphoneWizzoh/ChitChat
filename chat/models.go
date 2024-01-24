package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct{
	gorm.Model
	Username string
	Password string
}

func InitDatabase(){
	// Connect to SQLite database

	var err error
	db, err = gorm.Open(sqlite.Open("../chitchat.sqlite3"))
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate model structs
	db.AutoMigrate(&User{}) 
}