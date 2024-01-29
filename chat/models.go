package main

import (
	"errors"
	"log"

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

func (user *User) SetPassword(password string) error{
	hashedPassword, err := hashString(password)
	if err != nil{
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) CheckPassword(password string) error{
	inputHash, err := hashString(password)
	if err != nil{
		log.Println("Error in hashing during password verification")
		return err
	}

	if user.Password != inputHash{
		return errors.New("passwords do not match!")
	}
	return nil
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
		hashPassword, err := hashString(password)
		if err != nil{
			registrationError = err
		}
		db.Create(&User{
			Username: username,
			Email: email,
			Password: hashPassword,
		})
	}else{
		registrationError = errors.New("User already exists")
	}

	return registrationError
}

func GetUser(userId string)(User, bool){
	user, exists:= User{}, false
	return user, exists
}