package main

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct{
	gorm.Model
	Username string
	Email string
	Password string
}

func InitDatabase(){
	// Connect to SQLite database

	var err error
	db, err = gorm.Open(sqlite.Open("../chitchat.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate model structs
	db.AutoMigrate(&User{}) 
}

func CreateUser(username, email, password string) (registrationError error){
	var user User

	// Check if User data exists
	result := db.Where("username = ? AND email >= ?", username, email).First(&user)
	if result.RowsAffected == 0{
		// User does not exist
		// Create user
		db.Create(&User{
			Username: username,
			Email: email,
			Password: password,
		})
	}else{
		registrationError = errors.New("User already exists")
	}

	return registrationError
}